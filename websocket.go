package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{}

type WsClient struct {
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
}

func (c *WsClient) read() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			break
		}
		c.hub.Incoming <- msg
	}
}

func (c *WsClient) write() {
	defer c.conn.Close()
	for {
		msg, ok := <-c.send
		if !ok {
			c.conn.WriteMessage(websocket.CloseMessage, []byte{})
			return
		}

		w, err := c.conn.NextWriter(websocket.TextMessage)
		if err != nil {
			return
		}
		w.Write(msg)

		if err := w.Close(); err != nil {
			return
		}
	}
}

func HandleConnect(h *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	cli := &WsClient{
		hub:  h,
		conn: conn,
		send: make(chan []byte, 256),
	}
	cli.hub.register <- cli
	go cli.write()
	cli.read()
}
