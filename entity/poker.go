package entity

type RunInfo struct {
	DefaultDeal int
	Ante        Ante
	Deck        Deck
	PokerHands  *PokerHands
	Hands       int
	Discards    int
}

func NewRunInfo() *RunInfo {
	ante := NewAnte()
	ante.Base = ante.GetAnteBase()

	return &RunInfo{
		DefaultDeal: 8,
		Ante:        *ante,
		Deck:        NewDeck(),
		PokerHands:  NewPokerHands(),
		Hands:       4,
		Discards:    3,
	}
}

type Ante struct {
	Number int
	Base   int
	Blinds []Blind
}

type Blind struct {
	Name  string
	Multi float64
}

func NewAnte() *Ante {
	return &Ante{
		Number: 1,
		Blinds: []Blind{
			{
				Name:  "Small Blind",
				Multi: 1.0,
			},
			{
				Name:  "Big Blind",
				Multi: 1.5,
			},
			{
				Name:  "Final Blind",
				Multi: 2.0,
			},
		},
	}
}

func (a *Ante) GetAnteBase() int {
	switch a.Number {
	case 1:
		return 300
	case 2:
		return 800
	case 3:
		return 2800
	case 4:
		return 6000
	case 5:
		return 11000
	case 6:
		return 20000
	case 7:
		return 35000
	case 8:
		return 50000
	case 9:
		return 110000
	case 10:
		return 560000
	case 11:
		return 7200000
	case 12:
		return 300000000
	case 13:
		return 47000000000
	case 14:
		// return 2900e13
		return 2900 * 100000000000
	case 15:
		// return 7700e16
		return 7700 * 1000000000000
	case 16:
		// FIXME: return 8600e20
		return 8600
	default:
		return 0
	}
}
