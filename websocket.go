package main

import (
	"github.com/gorilla/websocket"
	"github.com/pulse0ne/gotunes/message"
	"github.com/satori/go.uuid"
	"net/http"
	"sync"
)

var wsUpgrader = websocket.Upgrader{}

type WsHub struct {
	listenMx    sync.RWMutex
	listeners   map[string]func(*message.WsMessage)
	connections map[string]*WsConnection
	register    chan *WsConnection
	unregister  chan *WsConnection
	Broadcast   chan *message.WsMessage
	Incoming    chan *message.WsMessage
	Outgoing    chan *message.WsMessage
}

func NewWsHub() *WsHub {
	h := &WsHub{
		listeners:   make(map[string]func(*message.WsMessage)),
		connections: make(map[string]*WsConnection),
		register:    make(chan *WsConnection),
		unregister:  make(chan *WsConnection),
		Broadcast:   make(chan *message.WsMessage, 1),
		Incoming:    make(chan *message.WsMessage),
		Outgoing:    make(chan *message.WsMessage, 1),
	}
	go h.run()
	return h
}

func (h *WsHub) AddListener(id string, f func(*message.WsMessage)) {
	h.listenMx.RLock()
	if _, ok := h.listeners[id]; !ok {
		h.listenMx.RUnlock()
		h.listenMx.Lock()
		h.listeners[id] = f
		h.listenMx.Unlock()
	} else {
		h.listenMx.RUnlock()
	}
}

func (h *WsHub) RemoveListener(id string) {
	h.listenMx.RLock()
	if _, ok := h.listeners[id]; ok {
		h.listenMx.RUnlock()
		h.listenMx.Lock()
		delete(h.listeners, id)
		h.listenMx.Unlock()
	} else {
		h.listenMx.RUnlock()
	}
}

func (h *WsHub) run() {
	for {
		select {
		case conn := <-h.register:
			LOG.Debug("registering connection")
			h.connections[conn.Id] = conn
		case conn := <-h.unregister:
			LOG.Debug("unregistering connection")
			if _, ok := h.connections[conn.Id]; ok {
				delete(h.connections, conn.Id)
				close(conn.Outbound)
				close(conn.Inbound)
			}
		case msg := <-h.Broadcast:
			for id, client := range h.connections {
				select {
				case client.Outbound <- msg:
				default:
					close(client.Outbound)
					close(client.Inbound)
					delete(h.connections, id)
				}
			}
		case msg := <-h.Outgoing:
			LOG.Debug("outgoing message")
			id := msg.ClientId
			if cli, ok := h.connections[id]; ok {
				cli.Outbound <- msg
			}
		case msg := <-h.Incoming:
			LOG.Debug("incoming message")
			h.listenMx.RLock()
			for _, listener := range h.listeners {
				listener(msg)
			}
			h.listenMx.RUnlock()
		}
	}
}

func HandleConnection(h *WsHub, w http.ResponseWriter, r *http.Request, mp func() *message.WsMessage) {
	c, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		LOG.Error("Could not upgrade connection")
		return
	}
	conn := NewWsConnection(h, c)
	conn.hub.register <- conn
	go conn.doWrite()
	conn.Outbound <- mp()
	conn.doRead()
}

//=====================================

type WsConnection struct {
	conn     *websocket.Conn
	hub      *WsHub
	Id       string
	Outbound chan *message.WsMessage
	Inbound  chan *message.WsMessage
}

func NewWsConnection(h *WsHub, c *websocket.Conn) *WsConnection {
	return &WsConnection{
		conn:     c,
		hub:      h,
		Id:       uuid.NewV4().String(),
		Outbound: make(chan *message.WsMessage),
		Inbound:  make(chan *message.WsMessage),
	}
}

func (c *WsConnection) doRead() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(1024)

	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				LOG.Error(err)
			}
			break
		}
		b, err := message.FromJsonBytes(msg)
		if err != nil {
			LOG.Error("Could not decode message:", err)
		} else {
			b.ClientId = c.Id
			c.hub.Incoming <- b
		}
	}
}

func (c *WsConnection) doWrite() {
	defer c.conn.Close()
	for {
		msg, ok := <-c.Outbound
		if !ok {
			c.conn.WriteMessage(websocket.CloseMessage, []byte{})
			return
		}

		w, err := c.conn.NextWriter(websocket.TextMessage)
		if err != nil {
			return
		}
		b, _ := message.ToJsonBytes(msg)
		w.Write(b)

		if err := w.Close(); err != nil {
			return
		}
	}
}
