package types

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"time"
)

type UTCTime struct {
	time.Time
}

var (
	ErrZeroValueTime = errors.New("Given time cannot be zero / null")
)

// NewUTCTime create a new UTCtime instance given a time.Time
func NewUTCTime(t time.Time) UTCTime {
	return UTCTime{Time: t}
}

// valuer / scanner interfaces

// Value is called when a value of this type is writtent o db
func (t UTCTime) Value() (driver.Value, error) {
	if t.IsZero() {
		// if given time is a zero return a zero as there's a db constraint (not null)
		return nil, ErrZeroValueTime
	}
	return t.Time.UTC().Format(time.RFC3339), nil
}

// Scan is called to read from db
func (t *UTCTime) Scan(value any) error {
	if value == nil {
		t.Time = time.Time{}
		return nil
	}
	if scannedTime, ok := value.(time.Time); ok { // TODO
		t.Time = scannedTime
		return nil
	}
	return fmt.Errorf("unsupported type for UTCTime: %T", value)
}
