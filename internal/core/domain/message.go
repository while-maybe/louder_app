package domain

import (
	"time"
)

type Message string
type RandomNumber uint

type MsgWithTime struct {
	CurrentLocalTime time.Time
	Message          string
}
