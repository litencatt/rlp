package entity

type Ante struct {
	BaseScore int
}

type RunInfo struct {
	Hands           int
	Discards        int
	Deck            Deck
	PokerHandLevels []PokerHandLevel
}

type PokerHandLevel struct {
	HandType HandType
	Level    int
}
