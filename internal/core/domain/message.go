package domain

import (
	"time"
)

type Message string

type RandomNumber uint

type MsgWithTime struct {
	Timestamp time.Time
	Content   string
}
