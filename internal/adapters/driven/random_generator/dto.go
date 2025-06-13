package randomgenerator

// TODO where do I put a struct mapper? I forgot :/

type DiceRollResponse struct {
	Roll []uint `json:"roll"`
	Sum  uint   `json:"sum"`
}
