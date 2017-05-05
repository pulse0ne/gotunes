package main

import (
	"fmt"
	"github.com/fhs/gompd/mpd"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"time"
)

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
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		conn.Stop()
		conn.Close()
		os.Exit(0)
	}()

	line := ""
	line1 := ""

	attrs, err := conn.ListAllInfo("/")
	if err != nil {
		log.Fatalln(err)
	}

	for _, v := range attrs {
		fmt.Println(v["file"])
		if err := conn.Add(v["file"]); err != nil {
			log.Fatalln("Could not add a file")
		}
	}

	err2 := conn.Play(-1)
	if err2 != nil {
		log.Fatalln(err2)
	}

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
