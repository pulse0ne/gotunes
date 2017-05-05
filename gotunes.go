package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/pulse0ne/gompd/mpd"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"time"
)

func fileserve() {
	log.Fatalln(http.ListenAndServe(":9999", http.FileServer(http.Dir("./public"))))
}

func attachExitHook(hook func(ch chan os.Signal)) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, os.Kill)
	go hook(c)
}

var wsUpgrader = websocket.Upgrader{}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	c, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("could not upgrade connection:", err)
		return
	}

	defer c.Close()

	for {
		mt, msg, err := c.ReadMessage()
		if err != nil {
			log.Println("could not read message:", err)
			break
		}
		if mt == websocket.TextMessage {
			fmt.Printf("WS received: %s\n", msg)
		}
	}
}

func main() {
	var conn *mpd.Client

	conn, err := mpd.Dial("tcp", "localhost:6600")
	if err != nil {
		fmt.Println("mpd is not running, attempting to start it...")
		cmd := exec.Command("mpd")
		_, err := cmd.CombinedOutput()
		if err != nil {
			log.Fatalln(err)
		}
		var ferr error
		conn, ferr = mpd.Dial("tcp", "localhost:6600")
		if ferr != nil {
			log.Fatalln(ferr)
		}
		fmt.Println("mpd successfully started")
	}

	// hook for interrupt
	attachExitHook(func(c chan os.Signal) {
		<-c // wait for the signal
		conn.Stop()
		conn.Close()
		os.Exit(0)
	})

	// start the fileserver
	go fileserve()

	http.HandleFunc("/websocket", wsHandler)

	attrs, err := conn.ListAllInfo("/")
	if err != nil {
		log.Fatalln(err)
	}

	for _, v := range attrs {
		fmt.Println(v)
		if err := conn.Add(v["file"]); err != nil {
			log.Fatalln("Could not add a file")
		}
	}

	err2 := conn.Play(-1)
	if err2 != nil {
		log.Fatalln(err2)
	}

	line := ""
	line1 := ""

	for {
		status, err := conn.Status()
		if err != nil {
			log.Fatalln(err)
		}

		song, err := conn.CurrentSong()
		if err != nil {
			log.Fatalln(err)
		}

		if status["state"] == "play" {
			line1 = fmt.Sprintf("%s - %s", song["Artist"], song["Title"])
		} else {
			line1 = fmt.Sprintf("State: %s", status["state"])
		}

		fmt.Printf("%s of %s\n", status["elapsed"], song["Time"])

		if line != line1 {
			line = line1
			fmt.Println(line)
		}
		time.Sleep(1e9)
	}
}

// server/application -> mpd -> JACK -> (qjackctl) -> JAMin EQ -> (qjackctl) -> output device
