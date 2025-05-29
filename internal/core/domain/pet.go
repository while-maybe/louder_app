package domain

import "time"

type Pet struct {
	name   string
	kind   PetKind
	age    time.Time
	tricks []Trick
}

type PetKind uint

const (
	unknownPetKind PetKind = iota
	cat
	dog
)

type Trick uint

const (
	unknownTrick Trick = iota
	sit
	stay
	lieDown
	fetch
	jump
	spin
)
