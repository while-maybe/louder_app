package dbcommon

type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	ErrHydrateWithNil   = Error("error attempted to hydrate person without data")
	ErrSaveNilPerson    = Error("error cannot save nil person to DB model")
	ErrConvertNilPerson = Error("error convert nil person to DB model")
	ErrEmptyID          = Error("error given ID is empty")
	ErrInvalidID        = Error("error invalid person ID format")
	ErrNotFoundInDB     = Error("error cannot find this ID in DB")
	ErrDBQueryFailed    = Error("error query has failed")
	ErrNilDomainPerson  = Error("error conversion returned nil domain person without error")
	// TODO individual errors for each adapter type
	ErrSavedButNotInDB = Error("error SQLx/Bun person saved but can't find in DB")
	ErrConvertPerson   = Error("error cannot convert SQLx/Bun data to a person")
)
