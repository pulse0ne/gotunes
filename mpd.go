package main

import (
	"fmt"
	"github.com/pulse0ne/gompd/mpd"
	"os/exec"
)

type MpdClient struct {
	client *mpd.Client
	Out    chan Message
	In     chan Message
}

func NewMpdClient(addr string) (*MpdClient, error) {
	var conn *mpd.Client
	conn, err := mpd.Dial("tcp", addr)
	if err != nil {
		fmt.Println("mpd is not running, attempting to start it...")
		cmd := exec.Command("mpd")
		out, err := cmd.CombinedOutput()
		if err != nil {
			return nil, err
		}
		fmt.Println(string(out))
		conn, err = mpd.Dial("tcp", addr)
		if err != nil {
			return nil, err
		}
		fmt.Println("mpd successfully started")
	}

	c := &MpdClient{
		client: conn,
		Out:    make(chan Message, 64),
		In:     make(chan Message, 64),
	}
	go c.writer()
	go c.reader()

	// exit hook
	GetHook().Add(func() {
		c.client.Pause(true)
		c.client.Close()
	})

	return c, nil
}

func (c *MpdClient) reader() {
	// TODO
}

func (c *MpdClient) writer() {
	// TODO
}
