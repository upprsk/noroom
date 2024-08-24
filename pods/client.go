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

type messageRequest struct {
	Id     int             `json:"id"`
	Method string          `json:"method"`
	Args   json.RawMessage `json:"args"`
	// Args   []any  `json:"args"`
}

type messageResponseCommon struct {
	Id  int    `json:"id"`
	Err string `json:"err"`
}

type messageResponseListPods struct {
	messageResponseCommon
	Pods []string `json:"pods"`
}

type messageResponseWithId struct {
	messageResponseCommon
	PodId string `json:"podId"`
}

func newMessageErr(id int, err error) messageResponseCommon {
	if err == nil {
		return messageResponseCommon{Id: id}
	}

	return messageResponseCommon{
		Id:  id,
		Err: err.Error(),
	}
}

func newMessageListPods(id int, pods []string) messageResponseListPods {
	return messageResponseListPods{
		messageResponseCommon: messageResponseCommon{Id: id},
		Pods:                  pods,
	}
}

func newMessageWithId(id int, podId string) messageResponseWithId {
	return messageResponseWithId{
		messageResponseCommon: messageResponseCommon{Id: id},
		PodId:                 podId,
	}
}

func newMessageEmpty(id int) messageResponseCommon {
	return messageResponseCommon{Id: id}
}

type messageEvent struct {
	Name string `json:"name"`
	Body any    `json:"body"`
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
	case "attachToPod":
		return c.methodAttachToPod(ctx, msg)
	case "sendToPod":
		return c.methodSendToPod(ctx, msg)
	}

	return nil
}

func (c *Client) methodListPods(ctx context.Context, msg messageRequest) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()

	pods, err := c.central.ContainersList(ctx)
	if err != nil {
		return c.sendMessage(newMessageErr(msg.Id, err))
	}

	return c.sendMessage(newMessageListPods(msg.Id, pods))
}

func (c *Client) methodCreatePod(ctx context.Context, msg messageRequest) error {
	var args []string
	if err := json.Unmarshal(msg.Args, &args); err != nil {
		return err
	}

	if len(args) != 1 {
		return fmt.Errorf(
			"invalid number of arguments to createPod method, expected 1, got %v",
			len(msg.Args),
		)
	}

	id, err := c.central.ContainerCreate(ctx, args[0])
	if err != nil {
		return c.sendMessage(newMessageErr(msg.Id, err))
	}

	return c.sendMessage(newMessageWithId(msg.Id, id))
}

func (c *Client) methodAttachToPod(ctx context.Context, msg messageRequest) error {
	var args []string
	if err := json.Unmarshal(msg.Args, &args); err != nil {
		return err
	}

	if len(args) != 1 {
		return fmt.Errorf(
			"invalid number of arguments to createPod method, expected 1, got %v",
			len(msg.Args),
		)
	}

	err := c.central.ContainerAttach(ctx, args[0], c)
	if err != nil {
		return c.sendMessage(newMessageErr(msg.Id, err))
	}

	return c.sendMessage(newMessageWithId(msg.Id, args[0]))
}

func (c *Client) methodSendToPod(ctx context.Context, msg messageRequest) error {
	var args [][]byte
	if err := json.Unmarshal(msg.Args, &args); err != nil {
		return err
	}

	if len(args) != 1 {
		return fmt.Errorf(
			"invalid number of arguments to createPod method, expected 1, got %v",
			len(msg.Args),
		)
	}

	err := c.central.containerSend(c, args[0])
	if err != nil {
		return c.sendMessage(newMessageErr(msg.Id, err))
	}

	return c.sendMessage(newMessageEmpty(msg.Id))
}

func (c *Client) sendMessage(msg any) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	c.send <- data
	return nil
}

func (c *Client) SendEvent(msg messageEvent) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	c.send <- data
	return nil
}
