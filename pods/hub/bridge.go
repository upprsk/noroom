package hub

import (
	"io"
	"log"

	"github.com/docker/docker/api/types"
)

type Bridge struct {
	containerStream types.HijackedResponse
}

func (b *Bridge) Connect(stream io.ReadWriteCloser) {
	go attachedContainerWritePump(b.containerStream.Conn, stream)
	go attachedContainerReadPump(b.containerStream.Conn, stream)
}

func (b *Bridge) Close() {
	b.containerStream.Close()
}

func attachedContainerWritePump(cont io.ReadCloser, bridge io.WriteCloser) {
	defer bridge.Close()
	defer cont.Close()

	for {
		buf := make([]byte, 512)
		n, err := cont.Read(buf)
		if err != nil {
			log.Println("error reading from container:", err)
			return
		}

		if _, err := bridge.Write(buf[:n]); err != nil {
			log.Println("error writing to bridge:", err)
			return
		}
	}
}

func attachedContainerReadPump(cont io.WriteCloser, bridge io.ReadCloser) {
	defer bridge.Close()
	defer cont.Close()

	for {
		buf := make([]byte, 512)
		n, err := bridge.Read(buf)
		if err != nil {
			log.Println("error reading from bridge:", err)
			return
		}

		if _, err := cont.Write(buf[:n]); err != nil {
			log.Println("error writing to container:", err)
			return
		}
	}
}
