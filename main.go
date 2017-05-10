package main

import (
	"log"
	"net/http"
)

func main() {
	//conn, err := mpd.Dial("tcp", "localhost:6600")
	//if err != nil {
	//	fmt.Println("mpd is not running, attempting to start it...")
	//	cmd := exec.Command("mpd")
	//	_, err := cmd.CombinedOutput()
	//	if err != nil {
	//		log.Fatalln(err)
	//	}
	//	var ferr error
	//	conn, ferr = mpd.Dial("tcp", "localhost:6600")
	//	if ferr != nil {
	//		log.Fatalln(ferr)
	//	}
	//	fmt.Println("mpd successfully started")
	//}

	//// cleanup on exit/interrupt
	//GetHook().Add(func() {
	//	conn.Pause(true)
	//	conn.Close()
	//})

	// start the fileserver
	go log.Fatalln(http.ListenAndServe(":9999", http.FileServer(http.Dir("./public"))))

	// websocket handling
	hub := NewHub()
	http.HandleFunc("/websocket", func(w http.ResponseWriter, r *http.Request) {
		HandleConnect(hub, w, r)
	})

	mpdclient, err := NewMpdClient("localhost:6600")
	if err != nil {
		log.Fatalln("Could not create mpd client:", err)
	}

	// channel routing
	for {
		select {
		case msg := <-mpdclient.Out:
			log.Println(msg)
			// TODO
		case msg := <-hub.Incoming:
			log.Println(msg)
			// TODO
		}
	}

	//attrs, err := conn.ListAllInfo("/")
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//
	//for _, v := range attrs {
	//	fmt.Println(v)
	//	if err := conn.Add(v["file"]); err != nil {
	//		log.Fatalln("Could not add a file")
	//	}
	//}
	//
	//err2 := conn.Play(-1)
	//if err2 != nil {
	//	log.Fatalln(err2)
	//}

	// TODO: the below might be a good fit for a goroutine
	//line := ""
	//line1 := ""
	//
	//for {
	//	status, err := conn.Status()
	//	if err != nil {
	//		log.Fatalln(err)
	//	}
	//
	//	song, err := conn.CurrentSong()
	//	if err != nil {
	//		log.Fatalln(err)
	//	}
	//
	//	if status["state"] == "play" {
	//		line1 = fmt.Sprintf("%s - %s", song["Artist"], song["Title"])
	//	} else {
	//		line1 = fmt.Sprintf("State: %s", status["state"])
	//	}
	//
	//	fmt.Printf("%s of %s\n", status["elapsed"], song["Time"])
	//
	//	if line != line1 {
	//		line = line1
	//		fmt.Println(line)
	//	}
	//	time.Sleep(75e7)
	//}
}

// server/application -> mpd -> JACK -> (qjackctl) -> JAMin EQ -> (qjackctl) -> output device
