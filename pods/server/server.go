package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"noroom/pods/hub"

	"github.com/quic-go/quic-go"
)

type Server struct {
	hub *hub.Hub
}

func NewServer(hub *hub.Hub) *Server {
	return &Server{
		hub: hub,
	}
}

func (s *Server) Start(port int) error {
	udpConn, err := net.ListenUDP("udp4", &net.UDPAddr{Port: port})
	if err != nil {
		return fmt.Errorf("failed to listen UDP: %w", err)
	}

	tr := quic.Transport{
		Conn: udpConn,
	}

	defer tr.Close()

	ln, err := tr.Listen(generateTLSConfig(), nil)
	if err != nil {
		return fmt.Errorf("failed to listen on transport: %w", err)
	}

	for {
		conn, err := ln.Accept(context.Background())
		if err != nil {
			var streamErr *quic.StreamError
			if errors.As(err, &streamErr) {
				errorCode := streamErr.ErrorCode
				if errorCode == 1000 {
					// actually, no error, just closed
					return nil
				}
			}

			return fmt.Errorf("failed to accept: %w", err)
		}

		log.Println("got new connection:", conn.RemoteAddr())

		cli := newConnClient(conn, s.hub)
		go func() {
			if err := cli.handle(); err != nil {
				log.Println("client error:", err)
			}
		}()
	}
}
