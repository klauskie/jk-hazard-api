package entity

type Card struct {
	ID      int64
	Src     string
	IsBlack bool
}

func (card *Card) IsEmpty() bool {
	return *card == Card{}
}