package main

import (
	"github.com/pulse0ne/gotunes/message"
	"sync"
)

type PlayTime struct {
	Total   int `json:"total"`
	Current int `json:"current"`
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
}

func (a *Track) Equal(b *Track) bool {
	return &a == &b || (a.Artist == b.Artist &&
		a.Album == b.Album &&
		a.Title == b.Title &&
		a.Duration == b.Duration &&
		a.TrackNum == b.TrackNum)

}

type Context struct {
	CtxType message.ContextType `json:"type"`
	Data    string              `json:"data"`
}

func (c *Context) Equal(b *Context) bool {
	return &c == &b || (c.CtxType == b.CtxType && c.Data == b.Data)
}

type Info struct {
	Time      PlayTime           `json:"time"`
	Playstate message.PlayState  `json:"playstate"`
	Volume    int                `json:"volume"`
	Track     Track              `json:"track"`
	Repeat    message.RepeatMode `json:"repeat"`
	Shuffle   bool               `json:"shuffle"`
	Context   Context            `json:"context"`
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

func (n *NowPlaying) Equal(b *NowPlaying) bool {
	return &n == &b || (n.Info.Time.Equal(&b.Info.Time) &&
		n.Info.Track.Equal(&b.Info.Track) &&
		n.Info.Playstate == b.Info.Playstate &&
		n.Info.Volume == b.Info.Volume &&
		n.Info.Repeat == b.Info.Repeat &&
		n.Info.Shuffle == b.Info.Shuffle &&
		n.Info.Context.Equal(&b.Info.Context))
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

func (n *NowPlaying) SetTrackRepeat(r message.RepeatMode) {
	n.infoMx.Lock()
	n.Info.Repeat = r
	n.infoMx.Unlock()
}

func (n *NowPlaying) SetShuffle(b bool) {
	n.infoMx.Lock()
	n.Info.Shuffle = b
	n.infoMx.Unlock()
}

func (n *NowPlaying) SetContext(c Context) {
	n.infoMx.Lock()
	n.Info.Context = c
	n.infoMx.Unlock()
}
