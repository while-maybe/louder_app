package domain

import (
	"fmt"
	"time"
)

type Message string
type Number uint

type MsgWithTime struct {
	CurrentLocalTime time.Time
	Message          string
}

func NewMsgWithTime(msg fmt.Stringer) *MsgWithTime {
	var msgContents string
	if msg == nil {
		msgContents = "mo message"
	} else {
		msgContents = msg.String()
	}

	return &MsgWithTime{
		CurrentLocalTime: time.Now(),
		Message:          msgContents,
	}
}
