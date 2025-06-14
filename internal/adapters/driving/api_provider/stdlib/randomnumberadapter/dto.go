package randomnumberadapter

import "louder/internal/core/domain"

type DiceRollDTO struct {
	Roll []uint `json:"roll"`
	Sum  uint   `json:"sum"`
}

type RandomNumberResponse struct {
	RandomNumber domain.RandomNumber `json:"random_number,omitempty"`
}

type DiceRollResponse struct {
	DiceRoll DiceRollDTO `json:"diceroll"`
}
