package main

import (
	"github.com/pulse0ne/gotunes/message"
	"sync"
)

type PlayTime struct {
	Total   int `json:"total"`
	Current int `json:"current"`
}

func (a *PlayTime) Copy() PlayTime {
	return PlayTime{
		Total:   a.Total,
		Current: a.Current,
	}
}

func (a *PlayTime) Equal(b *PlayTime) bool {
	return &a == &b || (a.Total == b.Total && a.Current == b.Current)
}

type Track struct {
	Artist   string `json:"artist"`
	Album    string `json:"album"`
	Title    string `json:"title"`
	Duration int    `json:"duration"`
	TrackNum string `json:"tracknum"`
	File     string `json:"file"`
}

func (a *Track) Copy() Track {
	return Track{
		Artist:   a.Artist,
		Album:    a.Album,
		Title:    a.Title,
		Duration: a.Duration,
		TrackNum: a.TrackNum,
		File:     a.File,
	}
}

func (a *Track) Equal(b *Track) bool {
	return &a == &b || (a.Artist == b.Artist &&
		a.Album == b.Album &&
		a.Title == b.Title &&
		a.Duration == b.Duration &&
		a.TrackNum == b.TrackNum &&
		a.File == b.File)

}

type Info struct {
	Time      PlayTime           `json:"time"`
	Playstate message.PlayState  `json:"playstate"`
	Volume    int                `json:"volume"`
	Track     Track              `json:"track"`
	Repeat    message.RepeatMode `json:"repeat"`
	Shuffle   bool               `json:"shuffle"`
}

type NowPlaying struct {
	infoMx sync.RWMutex
	Info   *Info
}

func NewNowPlaying() *NowPlaying {
	return &NowPlaying{
		Info: &Info{
			Playstate: message.STOPPED,
			Repeat:    message.OFF,
		},
	}
}

func (n *NowPlaying) Copy() NowPlaying {
	return NowPlaying{
		Info: &Info{
			Time:      n.Info.Time.Copy(),
			Playstate: n.Info.Playstate,
			Volume:    n.Info.Volume,
			Track:     n.Info.Track.Copy(),
			Repeat:    n.Info.Repeat,
			Shuffle:   n.Info.Shuffle,
		},
	}
}

func (n *NowPlaying) Equal(b *NowPlaying) bool {
	return &n == &b || (n.Info.Time.Equal(&b.Info.Time) &&
		n.Info.Track.Equal(&b.Info.Track) &&
		n.Info.Playstate == b.Info.Playstate &&
		n.Info.Volume == b.Info.Volume &&
		n.Info.Repeat == b.Info.Repeat &&
		n.Info.Shuffle == b.Info.Shuffle)
}

func (n *NowPlaying) GetInfo() Info {
	n.infoMx.RLock()
	i := *n.Info
	n.infoMx.RUnlock()
	return i
}

func (n *NowPlaying) SetTime(t PlayTime) {
	n.infoMx.Lock()
	n.Info.Time = t
	n.infoMx.Unlock()
}

func (n *NowPlaying) SetTimeCurrent(t int) {
	n.infoMx.Lock()
	n.Info.Time.Current = t
	n.infoMx.Unlock()
}

func (n *NowPlaying) SetTimeTotal(t int) {
	n.infoMx.Lock()
	n.Info.Time.Total = t
	n.infoMx.Unlock()
}

func (n *NowPlaying) SetPlaystate(p message.PlayState) {
	n.infoMx.Lock()
	n.Info.Playstate = p
	n.infoMx.Unlock()
}

func (n *NowPlaying) SetVolume(v int) {
	if v < 0 {
		v = 0
	} else if v > 100 {
		v = 100
	}
	n.infoMx.Lock()
	n.Info.Volume = v
	n.infoMx.Unlock()
}

func (n *NowPlaying) SetTrack(t Track) {
	n.infoMx.Lock()
	n.Info.Track = t
	n.infoMx.Unlock()
}

func (n *NowPlaying) SetTrackArtist(a string) {
	n.infoMx.Lock()
	n.Info.Track.Artist = a
	n.infoMx.Unlock()
}

func (n *NowPlaying) SetTrackAlbum(a string) {
	n.infoMx.Lock()
	n.Info.Track.Album = a
	n.infoMx.Unlock()
}

func (n *NowPlaying) SetTrackTitle(t string) {
	n.infoMx.Lock()
	n.Info.Track.Title = t
	n.infoMx.Unlock()
}

func (n *NowPlaying) SetTrackDuration(d int) {
	n.infoMx.Lock()
	n.Info.Track.Duration = d
	n.infoMx.Unlock()
}

func (n *NowPlaying) SetTrackNum(t string) {
	n.infoMx.Lock()
	n.Info.Track.TrackNum = t
	n.infoMx.Unlock()
}

func (n *NowPlaying) SetTrackFile(f string) {
	n.infoMx.Lock()
	n.Info.Track.File = f
	n.infoMx.Unlock()
}

func (n *NowPlaying) SetRepeat(r message.RepeatMode) {
	n.infoMx.Lock()
	n.Info.Repeat = r
	n.infoMx.Unlock()
}

func (n *NowPlaying) SetShuffle(b bool) {
	n.infoMx.Lock()
	n.Info.Shuffle = b
	n.infoMx.Unlock()
}
