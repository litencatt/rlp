package entity

import (
	"math/rand"
	"time"
)

type Suit string
type Rank string

const (
	Club    Suit = "Club"
	Diamond      = "Diamond"
	Heart        = "Heart"
	Spade        = "Spade"
)

const (
	Two   Rank = "2"
	Three      = "3"
	Four       = "4"
	Five       = "5"
	Six        = "6"
	Seven      = "7"
	Eight      = "8"
	Nine       = "9"
	Ten        = "T"
	Jack       = "J"
	Queen      = "Q"
	King       = "K"
	Ace        = "A"
)

type Trump struct {
	Suit Suit
	Rank Rank
}

func (t Trump) GetRankNumber() int {
	switch t.Rank {
	case Two:
		return 2
	case Three:
		return 3
	case Four:
		return 4
	case Five:
		return 5
	case Six:
		return 6
	case Seven:
		return 7
	case Eight:
		return 8
	case Nine:
		return 9
	case Ten:
	case Jack:
	case Queen:
	case King:
		return 10
	case Ace:
		return 11
	}
	return 0
}

type Deck []Trump

func NewDeck() Deck {
	deck := make(Deck, 0)
	suits := []Suit{Club, Diamond, Heart, Spade}
	ranks := []Rank{Two, Three, Four, Five, Six, Seven, Eight, Nine, Ten, Jack, Queen, King, Ace}
	for _, suit := range suits {
		for _, rank := range ranks {
			deck = append(deck, Trump{Suit: suit, Rank: rank})
		}
	}
	return deck
}

func (d Deck) Shuffle() {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(d), func(i, j int) {
		d[i], d[j] = d[j], d[i]
	})
}

func (d *Deck) Deal(n int) []Trump {
	hand := (*d)[:n]
	*d = (*d)[n:]
	return hand
}
