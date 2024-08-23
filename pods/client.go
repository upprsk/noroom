package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

type Client struct {
	central *Central
	conn    *websocket.Conn
	send    chan []byte

	// if the client has sent the initial handshake message
	opened bool
}

func NewClient(central *Central, conn *websocket.Conn) *Client {
	return &Client{
		central: central,
		conn:    conn,
		send:    make(chan []byte, 256),
	}
}

func (c *Client) readPump() {
	defer func() {
		c.conn.Close()
		c.central.Unregister(c)
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		if err := c.parseMessage(message); err != nil {
			log.Println("client error:", err)
			return // close
		}
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The central closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

type messageRequest struct {
	Id     int    `json:"id"`
	Method string `json:"method"`
	Args   []any  `json:"args"`
}

type messageResponse struct {
	Id   int    `json:"id"`
	Body any    `json:"body"`
	Err  string `json:"err"`
}

func (c *Client) parseMessage(data []byte) error {

	var msg messageRequest
	fmt.Println("data:", string(data), "msg", msg)

	if err := json.Unmarshal(data, &msg); err != nil {
		return err
	}

	if msg.Method == "open" {
		if c.opened {
			return fmt.Errorf("invalid open message received, client already open")
		}

		c.opened = true
		c.central.Register(c)
		return nil
	}

	if msg.Id == 0 {
		return fmt.Errorf("invalid id received: %d", msg.Id)
	}

	if !c.opened {
		return fmt.Errorf("invalid non-open message received, client is not open")
	}

	ctx := context.Background()
	switch msg.Method {
	case "listPods":
		return c.methodListPods(ctx, msg)
	case "createPod":
		return c.methodCreatePod(ctx, msg)
	}

	return nil
}

func (c *Client) methodListPods(ctx context.Context, msg messageRequest) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()

	pods, err := c.central.ContainersList(ctx)
	if err != nil {
		return c.sendMessage(messageResponse{
			Id:  msg.Id,
			Err: err.Error(),
		})
	}

	return c.sendMessage(messageResponse{
		Id:   msg.Id,
		Body: pods,
	})
}

func (c *Client) methodCreatePod(ctx context.Context, msg messageRequest) error {
	if len(msg.Args) != 1 {
		return fmt.Errorf(
			"invalid number of arguments to createPod method, expected 1, got %v",
			len(msg.Args),
		)
	}

	name, ok := msg.Args[0].(string)
	if !ok {
		return fmt.Errorf(
			"invalid arguments to createPod method, expected string at [0], got %v",
			msg.Args[0],
		)
	}

	fmt.Println("name:", name)

	id, err := c.central.ContainerCreate(ctx)
	if err != nil {
		return err
	}

	return c.sendMessage(messageResponse{
		Id:   msg.Id,
		Body: map[string]string{"id": id},
	})
}

func (c *Client) sendMessage(msg messageResponse) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	c.send <- data
	return nil
}
