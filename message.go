package main

import "encoding/json"

//==================================
// MessageType enum
type MessageType int

const (
	NOW_PLAYING MessageType = iota
	VIEW_UPDATE
	COMMAND
)

//==================================
// PlayState enum
type PlayState int

const (
	STOPPED PlayState = iota
	PLAYING
	PAUSED
)

//==================================
// Command enum
type Command int

const (
	SET_PLAYSTATE Command = iota
	SEEK_TO
	PLAY_NEXT
	PLAY_PREV
	SET_VOLUME
	SET_CONTEXT
	REQUEST_VIEW
	NEW_PLAYLIST
	SAVE_PLAYLIST
	ADD_TO_PLAYLIST
)

//==================================
// Repeat enum
type RepeatMode int

const (
	OFF RepeatMode = iota
	ALL
	ONE
)

//==================================
// Message contains all messages passed between client and server
type Message struct {
	MType   MessageType
	Payload interface{}
}

func NewMessage(m MessageType) *Message {
	return &Message{
		MType: m,
	}
}

func (m *Message) ToJsonBytes() ([]byte, error) {
	b, err := json.Marshal(m)
	if err != nil {
		return []byte{}, err
	}
	return b, nil
}

func FromJsonBytes(b []byte) (Message, error) {
	m := Message{}
	err := json.Unmarshal(b, m)
	if err != nil {
		return m, err
	}
	return m, nil
}
