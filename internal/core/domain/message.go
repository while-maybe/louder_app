package domain

import (
	"time"
)

type Message string

type MsgWithTime struct {
	Timestamp time.Time
	Content   string
}
