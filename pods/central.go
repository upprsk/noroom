package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type Central struct {
	docker *client.Client

	// registered clients
	clients map[*Client]*clientData

	// register events from clients
	register chan *Client

	// unregister events from clients
	unregister chan *Client
}

func NewCentral(ctx context.Context) (*Central, error) {
	docker, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	return &Central{
		docker:     docker,
		clients:    map[*Client]*clientData{},
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}, nil
}

func (c *Central) Close() {
	c.docker.Close()
}

func (c *Central) Start() {
	for {
		select {
		case client := <-c.register:
			c.clients[client] = &clientData{}
		case client := <-c.unregister:
			if data, ok := c.clients[client]; ok {
				delete(c.clients, client)
				close(client.send)

				data.Close()
			}
		}
	}
}

func (c *Central) Register(cli *Client) {
	log.Println("register:", cli)
	c.register <- cli
}

func (c *Central) Unregister(cli *Client) {
	log.Println("unregister:", cli)
	c.unregister <- cli
}

func (c *Central) ContainersList(ctx context.Context) ([]string, error) {
	containers, err := c.docker.ContainerList(ctx, container.ListOptions{})
	if err != nil {
		return nil, err
	}

	var ids []string
	for _, container := range containers {
		ids = append(ids, container.ID)
	}

	return ids, nil
}

func (c *Central) ContainerCreate(ctx context.Context, name string) (string, error) {
	resp, err := c.docker.ContainerCreate(ctx, &container.Config{
		Image:        "alpine",
		Cmd:          []string{"sh"},
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		OpenStdin:    true,
		Tty:          true,
		Env:          []string{"TERM=dumb"},
	}, nil, nil, nil, name)
	if err != nil {
		return "", err
	}

	if err := c.ContainerStart(ctx, resp.ID); err != nil {
		return resp.ID, err
	}

	if err := c.containerInit(ctx, resp.ID); err != nil {
		return resp.ID, err
	}

	return resp.ID, nil
}

func (c *Central) ContainerStart(ctx context.Context, id string) error {
	return c.docker.ContainerStart(ctx, id, container.StartOptions{})
}

func (c *Central) ContainerAttach(ctx context.Context, id string, cli *Client) error {
	data, ok := c.clients[cli]
	if !ok {
		return fmt.Errorf("tried to attach with client that is not registered")
	}

	if data.stream.Conn != nil {
		return fmt.Errorf("already attached")
	}

	resp, err := c.docker.ContainerAttach(ctx, id, container.AttachOptions{
		Stream: true,
		Stdin:  true,
		Stdout: true,
		Stderr: true,
	})
	if err != nil {
		return fmt.Errorf("failed to attach: %w", err)
	}

	data.stream = resp
	data.send = make(chan []byte, 256)

	go data.readPump(cli, id)
	go data.writePump()

	return nil
}

func (c *Central) containerSend(cli *Client, data []byte) error {
	cliData, ok := c.clients[cli]
	if !ok {
		return fmt.Errorf("tried to send with client that is not registered")
	}

	if cliData.stream.Conn == nil {
		return fmt.Errorf("tried to send to container that is not running")
	}

	cliData.Send(data)

	return nil
}

func (c *Central) containerInit(ctx context.Context, id string) error {
	resp, err := c.docker.ContainerAttach(ctx, id, container.AttachOptions{
		Stream: true,
		Stdin:  true,
		Stdout: true,
		Stderr: true,
	})
	if err != nil {
		return err
	}
	defer resp.Close()

	if _, err := fmt.Fprintf(resp.Conn, "apk add --update alpine-sdk\n"); err != nil {
		return err
	}

	s := bufio.NewScanner(resp.Conn)
	for s.Scan() {
		txt := s.Text()

		fmt.Println("\033[90m", txt, "\033[0m")

		if strings.HasPrefix(txt, "OK: ") {
			break
		}
	}

	fmt.Println("done")

	return nil
}

type clientData struct {
	stream types.HijackedResponse
	send   chan []byte
}

func (c *clientData) Send(data []byte) {
	c.send <- data
}

func (c *clientData) Close() {
	if c.stream.Conn != nil {
		c.stream.Close()
	}
}

func (c *clientData) readPump(cli *Client, id string) {
	defer func() {
		c.stream.Close()
		close(c.send)
		c.send = nil
	}()

	for {
		buf := make([]byte, 1024)

		n, err := c.stream.Conn.Read(buf)
		if err != nil {
			log.Println("error reading from container stream:", err)

			return
		}

		data := buf[:n]
		cli.SendEvent(messageEvent{
			Name: "podOut",
			Body: struct {
				Id   string `json:"id"`
				Data []byte `json:"data"`
			}{
				Id:   id,
				Data: data,
			},
		})
	}
}

func (c *clientData) writePump() {
	defer c.stream.Close()

	for data := range c.send {
		if _, err := c.stream.Conn.Write(data); err != nil {
			log.Println("error writing to container stream:", err)
			return
		}
	}
}
