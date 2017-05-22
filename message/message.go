package message

import "encoding/json"

//==================================
// MessageType enum
type MessageType int

const (
	_ MessageType = iota
	NOW_PLAYING
	VIEW_UPDATE
	COMMAND
)

//==================================
// PlayState enum
type PlayState int

const (
	_ PlayState = iota
	STOPPED
	PLAYING
	PAUSED
)

//==================================
// Command enum
type Command int

const (
	_ Command = iota
	SET_PLAYSTATE
	SEEK_TO
	PLAY_NEXT
	PLAY_PREV
	PLAY_FROM_CONTEXT
	SET_VOLUME
	SET_CONTEXT
	REQUEST_VIEW
	NEW_PLAYLIST
	SAVE_PLAYLIST
	ADD_TO_PLAYLIST
	SET_SHUFFLE
	SET_REPEAT_MODE
)

//==================================
// ContextType enum
type ContextType int

const (
	_ ContextType = iota
	ALL_ARTISTS
	ARTIST_DETAIL
	ALL_ALBUMS
	ALBUM_DETAIL
	ALL_TRACKS
	PLAYLIST
	PLAYLIST_DETAIL
)

//==================================
// Repeat enum
type RepeatMode int

const (
	_ RepeatMode = iota
	OFF
	ALL
	ONE
)

//==================================
// WsMessage contains all messages passed between client and server
type WsMessage struct {
	ClientId string      `json:"-"`
	MType    MessageType `json:"type"`
	Payload  interface{} `json:"payload"`
}

func ToJsonBytes(m *WsMessage) ([]byte, error) {
	b, err := json.Marshal(m)
	if err != nil {
		return []byte{}, err
	}
	return b, nil
}

func FromJsonBytes(b []byte) (*WsMessage, error) {
	m := &WsMessage{}
	err := json.Unmarshal(b, m)
	if err != nil {
		return m, err
	}
	return m, nil
}
