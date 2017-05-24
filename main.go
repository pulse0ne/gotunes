package main

import (
	"encoding/json"
	"errors"
	"github.com/pulse0ne/gompd/mpd"
	"github.com/pulse0ne/gotunes/logger"
	"github.com/pulse0ne/gotunes/message"
	"net/http"
	"os/exec"
	"strconv"
	"time"
)

var LOG = logger.NewLogger(logger.Ldebug)
var hub *WsHub
var nowplaying = NewNowPlaying()
var conn *mpd.Client

//================
func mpdConnect(addr string) *mpd.Client {
	var err error
	conn, err = mpd.Dial("tcp", addr)
	if err != nil {
		LOG.Info("mpd is not running, attempting to start it....")

		cmd := exec.Command("mpd")
		_, err := cmd.CombinedOutput()
		if err != nil {
			LOG.Fatal(err)
		}
		var ferr error
		conn, ferr = mpd.Dial("tcp", addr)
		if ferr != nil {
			LOG.Fatal(ferr)
		}
		LOG.Info("mpd successfully started")
	}
	LOG.Info("successfully connected to mpd")

	return conn
}

func mpdReport(m *mpd.Client, h *WsHub) {
	for {
		time.Sleep(750 * time.Millisecond)
		select {
		case <-m.Closed:
			return
		default:
			attr, err1 := m.CurrentSong()
			stat, err2 := m.Status()
			if err1 != nil {
				LOG.Error(err1)
				continue
			}
			if err2 != nil {
				LOG.Error(err2)
				continue
			}

			v, _ := strconv.Atoi(stat["volume"])
			nowplaying.SetVolume(v)
			s, _ := strconv.Atoi(stat["single"])
			if s == 1 {
				nowplaying.SetRepeat(message.ONE)
			} else {
				r, _ := strconv.Atoi(stat["repeat"])
				if r == 1 {
					nowplaying.SetRepeat(message.ALL)
				} else {
					nowplaying.SetRepeat(message.OFF)
				}
			}
			r, _ := strconv.Atoi(stat["random"])
			nowplaying.SetShuffle(r == 1)
			p := stat["state"]
			if p == "pause" {
				nowplaying.SetPlaystate(message.PAUSED)
			} else if p == "play" {
				nowplaying.SetPlaystate(message.PLAYING)
			} else {
				nowplaying.SetPlaystate(message.STOPPED)
			}

			if len(attr) > 0 {
				nowplaying.SetTrackArtist(attr["Artist"])
				nowplaying.SetTrackTitle(attr["Title"])
				nowplaying.SetTrackAlbum(attr["Album"])
				nowplaying.SetTrackNum(attr["Track"])
				nowplaying.SetTrackFile(attr["file"])

				t, err := strconv.Atoi(attr["Time"])
				if err != nil {
					nowplaying.SetTrackDuration(0)
					nowplaying.SetTimeTotal(0)
				} else {
					nowplaying.SetTrackDuration(t)
					nowplaying.SetTimeTotal(t)
				}
			}

			t, err := strconv.ParseFloat(stat["elapsed"], 32)
			if err != nil {
				nowplaying.SetTimeCurrent(0)
			} else {
				nowplaying.SetTimeCurrent(int(t))
			}

			h.Broadcast <- &message.WsMessage{
				MType:   message.NOW_PLAYING,
				Payload: nowplaying.GetInfo(),
			}
		}
	}
}

//==============
// TODO: imnprove error messages
func handleCommand(msg *message.WsMessage) error {
	payload, ok := msg.Payload.(map[string]interface{})
	if !ok {
		return errors.New("Could not decode command")
	}
	cmd, ok := payload["command"]
	if !ok {
		return errors.New("No command field found in payload")
	}
	command, ok := cmd.(float64)
	if !ok {
		return errors.New("Command type in unexpected format")
	}
	switch message.Command(command) {
	case message.SET_PLAYSTATE:
		LOG.Debug("SET_PLAYSTATE")
		data, ok := payload["data"]
		if !ok {
			return errors.New("No playstate provided")
		}
		ps, ok := data.(float64)
		if !ok {
			return errors.New("Data for playstate in unexpected format")
		}
		switch message.PlayState(ps) {
		case message.STOPPED:
			return conn.Stop()
		case message.PLAYING:
			if nowplaying.GetInfo().Playstate == message.PAUSED {
				return conn.Pause(false)
			} else {
				return conn.Play(-1)
			}
		case message.PAUSED:
			return conn.Pause(true)
		default:
			return errors.New("Unrecognized playstate")
		}
	case message.SEEK_TO:
		LOG.Debug("SEEK_TO")
		data, ok := payload["data"]
		if !ok {
			return errors.New("No seek position provided")
		}
		st, ok := data.(float64)
		if !ok {
			return errors.New("Seek position in unrecognized format")
		}
		return conn.SeekCur(int(st))
	case message.PLAY_NEXT:
		LOG.Debug("PLAY_NEXT")
		return conn.Next()
	case message.PLAY_PREV:
		LOG.Debug("PLAY_PREV")
		return conn.Previous()
	case message.PLAY_QUEUE_FROM_POSITION:
		LOG.Debug("PLAY_QUEUE_FROM_POSITION")
		data, ok := payload["data"]
		if !ok {
			return errors.New("No data provided for PLAY_QUEUE_FROM_POSITION command")
		}
		pos, ok := data.(float64)
		if !ok {
			return errors.New("Position in non-number format")
		}
		return conn.Play(int(pos))
	case message.SET_VOLUME:
		LOG.Debug("SET_VOLUME")
		data, ok := payload["data"]
		if !ok {
			return errors.New("No data provided for SET_VOLUME command")
		}
		v, ok := data.(float64)
		if !ok {
			return errors.New("Volume in non-number format")
		}
		return conn.SetVolume(int(v))
	case message.SET_SHUFFLE:
		LOG.Debug("SET_SHUFFLE")
		data, ok := payload["data"]
		if !ok {
			return errors.New("No data provided for SET_SHUFFLE command")
		}
		shuffle, ok := data.(bool)
		if !ok {
			return errors.New("Shuffle is not in boolean format")
		}
		return conn.Random(shuffle)
	case message.SET_REPEAT_MODE:
		LOG.Debug("SET_REPEAT_MODE")
		data, ok := payload["data"]
		if !ok {
			return errors.New("No data provided for SET_REPEAT_MODE")
		}
		rm, ok := data.(float64)
		if !ok {
			return errors.New("Repeat mode in unrecognized format")
		}
		switch message.RepeatMode(rm) {
		case message.OFF:
			if nowplaying.GetInfo().Repeat == message.ONE {
				return conn.Single(false)
			} else {
				return conn.Repeat(false)
			}
		case message.ALL:
			return conn.Repeat(true)
		case message.ONE:
			return conn.Single(true)
		default:
			return errors.New("Unrecognized repeat mode")
		}
	case message.REQUEST_VIEW:
		LOG.Debug("REQUEST_VIEW")
		data, ok := payload["data"]
		if !ok {
			return errors.New("No data provided for view request")
		}
		view, ok := data.(float64)
		if !ok {
			return errors.New("Data for view request in unexpected format")
		}
		omsg := &message.WsMessage{
			ClientId: msg.ClientId,
			MType:    message.VIEW_UPDATE,
		}

		// TODO: refactor for less repetition
		switch message.ViewType(view) {
		case message.QUEUE:
			attr, err := conn.PlaylistInfo(-1, -1)
			if err != nil {
				return err
			}
			omsg.Payload = map[string]interface{}{
				"type": message.ViewType(view),
				"data": attr,
			}
		case message.ALL_ARTISTS:
			attr, err := conn.List("artist")
			if err != nil {
				return err
			}
			omsg.Payload = map[string]interface{}{
				"type": message.ViewType(view),
				"data": attr,
			}
		case message.ARTIST_DETAIL:
			artist, ok := payload["detail"]
			if !ok {
				return errors.New("No detail provided")
			}
			a, ok := artist.(string)
			if !ok {
				return errors.New("could not decode string")
			}
			attr, err := conn.Find("artist", a)
			if err != nil {
				return err
			}
			omsg.Payload = map[string]interface{}{
				"type": message.ViewType(view),
				"data": attr,
			}
		case message.ALL_ALBUMS:
			attr, err := conn.List("album")
			if err != nil {
				return err
			}
			omsg.Payload = map[string]interface{}{
				"type": message.ViewType(view),
				"data": attr,
			}
		case message.ALBUM_DETAIL:
			album, ok := payload["detail"]
			if !ok {
				return errors.New("No detail provided")
			}
			a, ok := album.(string)
			if !ok {
				return errors.New("could not decode string")
			}
			attr, err := conn.Find("album", a)
			if err != nil {
				return err
			}
			omsg.Payload = map[string]interface{}{
				"type": message.ViewType(view),
				"data": attr,
			}
		case message.ALL_TRACKS:
			attr, err := conn.ListAllInfo("/")
			if err != nil {
				return err
			}
			omsg.Payload = map[string]interface{}{
				"type": message.ViewType(view),
				"data": attr,
			}
		case message.PLAYLIST:
			attr, err := conn.ListPlaylists()
			if err != nil {
				return err
			}
			omsg.Payload = map[string]interface{}{
				"type": message.ViewType(view),
				"data": attr,
			}
		case message.PLAYLIST_DETAIL:
			playlist, ok := payload["detail"]
			if !ok {
				return errors.New("No detail provided")
			}
			pl, ok := playlist.(string)
			if !ok {
				return errors.New("Could not decode string")
			}
			attr, err := conn.PlaylistContents(pl)
			if err != nil {
				return err
			}
			omsg.Payload = map[string]interface{}{
				"type": message.Command(command),
				"data": attr,
			}
		default:
			return errors.New("Received an unsupported view type")
		}

		hub.Outgoing <- omsg
	case message.ADD_TO_QUEUE:
		LOG.Debug("ADD_TO_QUEUE")
		data, ok := payload["data"]
		if !ok {
			return errors.New("No data provided for ADD_TO_QUEUE command")
		}
		uri, ok := data.(string)
		if !ok {
			return errors.New("Could not decode uri string")
		}
		err := conn.Add(uri)
		if err != nil {
			return err
		}
	case message.SAVE_AS_PLAYLIST:
		LOG.Debug("SAVE_AS_PLAYLIST")
		// TODO
	case message.SAVE_PLAYLIST:
		LOG.Debug("SAVE_PLAYLIST")
		// TODO
	case message.DELETE_PLAYLIST:
		LOG.Debug("DELETE_PLAYLIST")
		// TODO
	case message.LOAD_PLAYLIST:
		LOG.Debug("LOAD_PLAYLIST")
		// TODO
	default:
		return errors.New("Unrecognized command")
	}
	return nil
}

func messageHandler(msg *message.WsMessage) {
	s, _ := json.Marshal(msg)
	LOG.Debug(string(s))
	switch msg.MType {
	case message.COMMAND:
		err := handleCommand(msg)
		if err != nil {
			LOG.Error(err)
		}
	default:
		LOG.Info("Got an unrecognized message type")
	}
}

//===============
func main() {
	// connect to mpd -- will exit fatally if connection cannot be made
	conn = mpdConnect("localhost:6600")
	defer func() {
		conn.Pause(true)
		conn.Close()
	}()

	hub = NewWsHub()
	hub.AddListener("main", messageHandler)

	http.Handle("/", http.FileServer(http.Dir("./public")))
	http.HandleFunc("/websocket", func(w http.ResponseWriter, r *http.Request) {
		LOG.Debug("New Websocket Connection")
		HandleConnection(hub, w, r, func() *message.WsMessage {
			return &message.WsMessage{
				MType:   message.NOW_PLAYING,
				Payload: nowplaying.GetInfo(),
			}
		})
	})

	go mpdReport(conn, hub)

	LOG.Fatal(http.ListenAndServe(":9999", nil))
}

// server/application -> mpd -> JACK -> (qjackctl) -> JAMin EQ -> (qjackctl) -> output device
