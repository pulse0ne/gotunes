package main

type Hub struct {
	clients    map[*WsClient]bool
	register   chan *WsClient
	unregister chan *WsClient
	Broadcast  chan []byte
	Incoming   chan []byte
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*WsClient]bool),
		register:   make(chan *WsClient, 1),
		unregister: make(chan *WsClient, 1),
		Broadcast:  make(chan []byte, 256),
		Incoming:   make(chan []byte, 256),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case cli := <-h.register:
			h.clients[cli] = true
		case cli := <-h.unregister:
			if _, ok := h.clients[cli]; ok {
				delete(h.clients, cli)
				close(cli.send)
			}
		case msg := <-h.Broadcast:
			for cli := range h.clients {
				select {
				case cli.send <- msg:
				default:
					close(cli.send)
					delete(h.clients, cli)
				}
			}
		}
	}
}
