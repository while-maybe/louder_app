package randomnumberadapter

import "louder/internal/core/domain"

func toRandomNumberDTO(p *domain.RandomDice) *DiceRollDTO {

	if p == nil {
		return nil
	}

	return &DiceRollDTO{
		Roll: p.Roll,
		Sum:  p.RollSum,
	}
}
