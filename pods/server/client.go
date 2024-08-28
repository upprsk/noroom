package server

import (
	"context"
	"log"
	"noroom/pods/hub"
	"noroom/rpc"
	"time"

	"github.com/quic-go/quic-go"
)

type connClient struct {
	conn quic.Connection
	hub  *hub.Hub
}

func newConnClient(conn quic.Connection, hub *hub.Hub) connClient {
	return connClient{
		conn: conn,
		hub:  hub,
	}
}

func (cc *connClient) close(code quic.ApplicationErrorCode, message string) error {
	return cc.conn.CloseWithError(code, message)
}

func (cc *connClient) handle() error {
	defer cc.close(1001, "unworthy")

	for {
		stream, err := cc.conn.AcceptStream(context.Background())
		if err != nil {
			return err
		}

		log.Println("got new stream:", stream.StreamID())
		cli := newStreamClient(stream, cc.hub)
		go func() {
			if err := cli.handle(); err != nil {
				log.Println("client stream error:", err)
			}
		}()
	}
}

// ============================================================================

const (
	defaultTimeout = time.Second * 10
)

type streamClient struct {
	stream quic.Stream
	rpc    *rpc.RpcServer
}

func newStreamClient(stream quic.Stream, hub *hub.Hub) streamClient {
	return streamClient{
		stream: stream,
		rpc:    rpc.NewRpcServer(stream, defaultTimeout, hub),
	}
}

func (sc *streamClient) close() error {
	if sc.stream != nil {
		return sc.stream.Close()
	}

	return nil
}

func (sc *streamClient) handle() error {
	defer sc.close()

	for {
		detach, err := sc.rpc.HandleOne(context.Background())
		if err != nil {
			return err
		}

		if detach {
			sc.stream = nil
			return nil
		}
	}
}

func (sc *streamClient) handleOne() (bool, error) {
	return sc.rpc.HandleOne(context.Background())
}
