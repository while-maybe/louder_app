package service

type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	ErrInvalidPersonData = Error("error invalid data received")
)
