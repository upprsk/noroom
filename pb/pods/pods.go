package pods

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"

	"noroom/rpc"

	"github.com/quic-go/quic-go"
)

type PodServerManager struct {
	podServers map[string]*podServer

	mutex sync.Mutex
}

func NewPodServerManager() *PodServerManager {
	return &PodServerManager{
		podServers: map[string]*podServer{},
		mutex:      sync.Mutex{},
	}
}

// in case the error is a connect error, the server is still added to the map
func (m *PodServerManager) Add(id, addr string) error {
	resolved, err := net.ResolveUDPAddr("udp4", addr)
	if err != nil {
		return fmt.Errorf("failed to resolve address %s: %w", addr, err)
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	if _, exists := m.podServers[id]; exists {
		return fmt.Errorf("server with id %v already added", id)
	}

	srv := newPodServer(resolved)
	go srv.start()

	m.podServers[id] = srv

	return srv.reconnect()
}

func (m *PodServerManager) Del(id string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if srv, ok := m.podServers[id]; ok {
		srv.close()
	}

	delete(m.podServers, id)

	return nil
}

func (m *PodServerManager) Update(id, addr string) error {
	if err := m.Del(id); err != nil {
		return err
	}

	return m.Add(id, addr)
}

func (m *PodServerManager) AddNewPodToServer(serverId, podName, podImage string) (string, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	podServer, ok := m.podServers[serverId]
	if !ok {
		return "", fmt.Errorf("no such server with id %v", serverId)
	}

	podId, err := podServer.createNewPod(podName, podImage)
	if err != nil {
		return "", fmt.Errorf("error adding new pod: %w", err)
	}

	return podId, nil
}

func (m *PodServerManager) AddExistingPodToServer(serverId, podId string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	podServer, ok := m.podServers[serverId]
	if !ok {
		return fmt.Errorf("no such server with id %v", serverId)
	}

	if err := podServer.addExistingPod(podId); err != nil {
		return fmt.Errorf("error adding existing pod: %w", err)
	}

	return nil
}

func (m *PodServerManager) AddExistingPodToServerWithoutConnect(serverId, podId string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	podServer, ok := m.podServers[serverId]
	if !ok {
		return fmt.Errorf("no such server with id %v", serverId)
	}

	if err := podServer.addExistingPodWithoutConnect(podId); err != nil {
		return fmt.Errorf("error adding existing pod: %w", err)
	}

	return nil
}

func (m *PodServerManager) DeletePodFromServer(serverId, podId string, timeout time.Duration) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	podServer, ok := m.podServers[serverId]
	if !ok {
		return fmt.Errorf("no such server with id %v", serverId)
	}

	if err := podServer.deletePod(podId, timeout); err != nil {
		return fmt.Errorf("error deleting pod: %w", err)
	}

	return nil
}

func (m *PodServerManager) StartPodById(podId string, timeout time.Duration) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// log.Println("StartPodById:", podId)
	srv, pod := m.findPodById(podId)
	if srv == nil {
		return fmt.Errorf("no such pod with id %v", podId)
	}

	if err := pod.start(timeout); err != nil {
		srv.reconnectIfNetErr(err)
		return err
	}

	return nil
}

func (m *PodServerManager) StopPodById(podId string, timeout time.Duration) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	srv, pod := m.findPodById(podId)
	if srv == nil {
		return fmt.Errorf("no such pod with id %v", podId)
	}

	// log.Println("StopPodById:", podId)
	if err := pod.stop(timeout); err != nil {
		srv.reconnectIfNetErr(err)
		return err
	}

	return nil
}

func (m *PodServerManager) KillPodById(podId string, timeout time.Duration) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	srv, pod := m.findPodById(podId)
	if srv == nil {
		return fmt.Errorf("no such pod with id %v", podId)
	}

	// log.Println("KillPodById:", podId)
	if err := pod.kill(timeout); err != nil {
		srv.reconnectIfNetErr(err)
		return err
	}

	return nil
}

func (m *PodServerManager) InspectPodById(podId string) (*rpc.ContainerInspectResult, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	srv, pod := m.findPodById(podId)
	if srv == nil {
		return nil, fmt.Errorf("no such pod with id %v", podId)
	}

	// log.Println("InspectPodById:", podId)
	data, err := pod.inspect()
	if err != nil {
		srv.reconnectIfNetErr(err)
		return nil, err
	}

	return data, nil
}

func (m *PodServerManager) AttachPodById(podId string) (quic.Stream, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	srv, pod := m.findPodById(podId)
	if srv == nil {
		return nil, fmt.Errorf("no such pod with id %v", podId)
	}

	// log.Println("InspectPodById:", podId)
	err := pod.attach()
	if err != nil {
		srv.reconnectIfNetErr(err)
		return nil, err
	}

	stream := pod.stream

	if err := srv.openRpcForPod(podId); err != nil {
		srv.reconnectIfNetErr(err)
		return nil, err
	}

	return stream, nil
}

func (m *PodServerManager) findPodById(podId string) (*podServer, *podInstance) {
	for _, srv := range m.podServers {
		pod, ok := srv.pods[podId]
		if ok {
			return srv, pod
		}
	}

	return nil, nil
}

// ============================================================================

type podServer struct {
	addr net.Addr
	tr   *quic.Transport
	conn quic.Connection

	pods map[string]*podInstance
	cmds chan podServerCmd
}

type podServerCmd struct {
	exec func() error
	ret  chan error
}

func newPodServer(addr net.Addr) *podServer {
	s := &podServer{
		addr: addr,
		pods: map[string]*podInstance{},
		cmds: make(chan podServerCmd),
	}

	return s
}

func (p *podServer) start() {
	for cmd := range p.cmds {
		cmd.ret <- cmd.exec()
	}
}

func (p *podServer) close() {
	p.execCmd(func() error {
		p.execCloseFailure()
		return nil
	})

	close(p.cmds)
}

func (p *podServer) reconnect() error {
	return p.execCmd(func() error {
		return p.execReconnect()
	})
}

func (p *podServer) reconnectIfNetErr(err error) {
	p.execCmd(func() error {
		p.execReconnectIfNet(err)
		return nil
	})
}

func (p *podServer) openRpcForPod(podId string) error {
	return p.execCmd(func() error {
		return p.execOpenRpcForPod(podId)
	})
}

func (p *podServer) createNewPod(name, image string) (string, error) {
	vid := new(string)

	if err := p.execCmd(func() error {
		id, err := p.execCreatePod(name, image)
		*vid = id

		return err
	}); err != nil {
		return "", err
	}

	return *vid, nil
}

func (p *podServer) addExistingPod(podId string) error {
	return p.execCmd(func() error {
		return p.execAddExistingPod(podId)
	})
}

func (p *podServer) addExistingPodWithoutConnect(podId string) error {
	return p.execCmd(func() error {
		return p.execAddExistingPodWithoutConnect(podId)
	})
}

func (p *podServer) updateConnectedPodUnlocked(podId string) error {
	return p.execCmd(func() error {
		return p.execUpdateConnectedPod(podId)
	})
}

func (p *podServer) deletePod(podId string, timeout time.Duration) error {
	return p.execCmd(func() error {
		return p.execDeletePod(podId, timeout)
	})
}

// ----------------------------------------------------------------------------

func (p *podServer) openStream() (quic.Stream, error) {
	// log.Println("podServer.openStream()")
	conn, err := p.getConnection()
	if err != nil {
		p.execCloseFailure()
		return nil, err
	}

	stream, err := conn.OpenStream()
	if err != nil {
		p.execCloseFailure()
		return nil, err
	}

	return stream, nil
}

func (p *podServer) getConnection() (quic.Connection, error) {
	if p.tr == nil {
		tr, err := makeQuicTransport()
		if err != nil {
			return nil, err
		}

		p.tr = tr
	}

	if p.conn == nil {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second) // 3s handshake timeout
		defer cancel()

		conn, err := connectQuic(ctx, p.tr, p.addr)
		if err != nil {
			return nil, err
		}

		p.conn = conn
	}

	return p.conn, nil
}

func (p *podServer) hasConnection() bool {
	return p.conn != nil
}

func connectQuic(ctx context.Context, tr *quic.Transport, addr net.Addr) (quic.Connection, error) {
	return tr.Dial(
		ctx,
		addr,
		&tls.Config{
			InsecureSkipVerify: true,
			NextProtos:         []string{"noroom-rpc"},
		},
		&quic.Config{
			KeepAlivePeriod: 30 * time.Second,
		},
	)
}

func makeQuicTransport() (*quic.Transport, error) {
	conn, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4zero, Port: 0})
	if err != nil {
		return nil, fmt.Errorf("failed to listen UDP for quic: %w", err)
	}

	return &quic.Transport{Conn: conn}, nil
}

// ----------------------------------------------------------------------------

func (p *podServer) execCmd(exec func() error) error {
	ret := make(chan error)

	p.cmds <- podServerCmd{
		exec: exec,
		ret:  ret,
	}

	return <-ret
}

func (p *podServer) execOpenRpcForPod(podId string) error {
	stream, err := p.openStream()
	if err != nil {
		return err
	}

	rpc := rpc.NewRpcClient(stream)
	p.pods[podId] = newPodInstance(podId, stream, rpc)
	// log.Println("opened new stream for pod with id:", podId)

	return nil
}

func (p *podServer) execCreatePod(name, image string) (string, error) {
	stream, err := p.openStream()
	if err != nil {
		return "", err
	}

	rpc := rpc.NewRpcClient(stream)

	podId, err := rpc.Create(name, image)
	if err != nil {
		p.execReconnectIfNet(err)
		return "", err
	}

	if _, exists := p.pods[podId]; exists {
		return "", fmt.Errorf("pod with id %v already added", podId)
	}

	p.pods[podId] = newPodInstance(podId, stream, rpc)
	// log.Println("created new pod with id:", podId)

	return podId, nil
}

func (p *podServer) execAddExistingPodWithoutConnect(podId string) error {
	if _, exists := p.pods[podId]; exists {
		return fmt.Errorf("pod with id %v already added", podId)
	}

	var stream quic.Stream
	if p.hasConnection() {
		s, err := p.openStream()
		if err != nil {
			log.Println("failed to open stream for pod with id", podId)
		} else {
			stream = s
		}
	}

	rpc := rpc.NewRpcClient(stream)
	p.pods[podId] = newPodInstance(podId, stream, rpc)
	// log.Println("added existing pod with id without connecting:", podId)

	return nil
}

func (p *podServer) execAddExistingPod(podId string) error {
	if _, exists := p.pods[podId]; exists {
		return fmt.Errorf("pod with id %v already added", podId)
	}

	stream, err := p.openStream()
	if err != nil {
		return err
	}

	rpc := rpc.NewRpcClient(stream)
	p.pods[podId] = newPodInstance(podId, stream, rpc)
	// log.Println("added existing pod with id:", podId)

	return nil
}

func (p *podServer) execDeletePod(podId string, timeout time.Duration) error {
	pod, exists := p.pods[podId]
	if !exists {
		return fmt.Errorf("pod with id %v does not exist", podId)
	}

	delete(p.pods, podId)
	if err := pod.kill(timeout); err != nil {
		log.Printf("failed to kill pod %v: %v", podId, err)
	}

	if err := pod.delete(); err != nil {
		log.Printf("failed to delete pod %v: %v", podId, err)
	}

	pod.close()

	return nil
}

func (p *podServer) execUpdateConnectedPod(podId string) error {
	pod, exists := p.pods[podId]
	if !exists {
		return fmt.Errorf("pod with id %v does not exist", podId)
	}

	pod.close()

	stream, err := p.openStream()
	if err != nil {
		return err
	}

	rpc := rpc.NewRpcClient(stream)
	p.pods[podId] = newPodInstance(podId, stream, rpc)
	// log.Println("added existing pod with id:", podId)

	return nil
}

func (p *podServer) execCloseFailure() {
	p.execClose(1001, "failure")
}

func (p *podServer) execCloseReconnect() {
	p.execClose(1002, "reconnect")
}

func (p *podServer) execReconnectIfNet(err error) {
	if nerr, ok := err.(net.Error); ok && nerr.Timeout() {
		// log.Println("got timeout error:", nerr)
		if err := p.execReconnect(); err != nil {
			log.Println("failure in execReconnectIfNet:", err)
		}
	} else if errors.Is(err, io.EOF) {
		// log.Println("got EOF error:", nerr)
		if err := p.execReconnect(); err != nil {
			log.Println("failure in execReconnectIfNet:", err)
		}
	} else if errors.Is(err, rpc.ErrNilStream) {
		// log.Println("got nil stream error:", nerr)
		if err := p.execReconnect(); err != nil {
			log.Println("failure in execReconnectIfNet:", err)
		}
	}
}

func (p *podServer) execClose(code quic.ApplicationErrorCode, message string) {
	// log.Printf("podServer.execClose(%v, %v)", code, message)

	if p.tr != nil {
		p.tr.Close()
		p.tr = nil
	}

	if p.conn != nil {
		p.conn.CloseWithError(code, message)
		p.conn = nil
	}
}

func (p *podServer) execConnect() error {
	// log.Printf("podServer.execConnect(%v)", p.addr)

	if p.tr != nil {
		return fmt.Errorf("pod server already as a connection transport")
	}

	if p.conn != nil {
		return fmt.Errorf("pod server already as a connection")
	}

	tr, err := makeQuicTransport()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second) // 3s handshake timeout
	defer cancel()

	conn, err := connectQuic(ctx, tr, p.addr)
	if err != nil {
		tr.Close()
		return err
	}

	p.tr = tr
	p.conn = conn

	// log.Println("connected")

	return nil
}

func (p *podServer) execReconnect() error {
	p.execCloseReconnect()
	if err := p.execConnect(); err != nil {
		return err
	}

	var errs []error
	for _, v := range p.pods {
		err := p.execUpdateConnectedPod(v.podId)
		if err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}

// ============================================================================

type podInstance struct {
	podId string

	stream quic.Stream
	rpc    *rpc.RpcClient
}

func newPodInstance(podId string, stream quic.Stream, rpc *rpc.RpcClient) *podInstance {
	return &podInstance{
		podId:  podId,
		stream: stream,
		rpc:    rpc,
	}
}

func (p *podInstance) close() {
	if p.stream != nil {
		p.stream.Close()
	}
}

func (p *podInstance) kill(timeout time.Duration) error {
	return p.rpc.Kill(p.podId, timeout)
}

func (p *podInstance) start(timeout time.Duration) error {
	return p.rpc.Start(p.podId, timeout)
}

func (p *podInstance) stop(timeout time.Duration) error {
	return p.rpc.Stop(p.podId, timeout)
}

func (p *podInstance) delete() error {
	return p.rpc.Delete(p.podId)
}

func (p *podInstance) inspect() (*rpc.ContainerInspectResult, error) {
	return p.rpc.Inspect(p.podId)
}

func (p *podInstance) attach() error {
	return p.rpc.Attach(p.podId)
}
