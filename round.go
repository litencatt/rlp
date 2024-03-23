package rlp

import (
	"strings"

	"github.com/litencatt/rlp/domain"
	"github.com/litencatt/rlp/entity"
)

type PokerRound struct {
	Deck          entity.Deck
	TotalScore    int
	HandCards     []entity.Trump
	RemainCards   []entity.Trump
	SelectedCards []entity.Trump
	Hands         int
	DisCards      int
	ScoreAtLeast  int
}

func NewRound() *PokerRound {
	deck := entity.NewDeck()
	deck.Shuffle()

	return &PokerRound{
		Deck:         deck,
		Hands:        4,
		DisCards:     3,
		ScoreAtLeast: 10,
	}
}

func (p *PokerRound) Start(defaultDeal int) []string {
	p.HandCards = p.Deck.Draw(defaultDeal)

	return p.HandCardString()
}

func (p *PokerRound) DrawCard(drawNum int) {
	p.HandCards = nil

	drawCards := p.Deck.Draw(drawNum)
	p.HandCards = append(p.HandCards, p.RemainCards...)
	p.HandCards = append(p.HandCards, drawCards...)

	domain.SortHand(p.HandCards)
}

func (p *PokerRound) HandCardString() []string {
	var cards []string
	for _, card := range p.HandCards {
		cards = append(cards, card.String())
	}
	return cards
}

func (p *PokerRound) RemainCardString() []string {
	var cards []string
	for _, card := range p.RemainCards {
		cards = append(cards, card.String())
	}
	return cards
}

func (p *PokerRound) SelectCards(cards []string) int {
	// Convert select cards to Trump entity
	var selectCards []entity.Trump
	for _, card := range cards {
		// extract rank and suit from card string
		rank := strings.Split(card, " of ")[0]
		suit := strings.Split(card, " of ")[1]
		// Find the card from hand
		for _, c := range p.HandCards {
			if string(c.Rank) == rank && string(c.Suit) == suit {
				selectCards = append(selectCards, c)
				break
			}
		}
	}
	p.SelectedCards = selectCards

	// Calc the RemainCards cards
	var RemainCardsCards []entity.Trump
	for _, card := range p.HandCards {
		if !domain.Contains(selectCards, card) {
			RemainCardsCards = append(RemainCardsCards, card)
		}
	}
	p.RemainCards = RemainCardsCards

	return len(p.SelectedCards)
}

func (p *PokerRound) PlayHand() (entity.HandType, int) {
	handType := entity.EvaluateHand(p.SelectedCards)
	score := entity.GetScore(handType)
	p.TotalScore += score

	return handType, score
}

func (p *PokerRound) GetTotalScore() int {
	return p.TotalScore
}

func (p *PokerRound) GetHands() int {
	return p.Hands
}

func (p *PokerRound) GetDisCards() int {
	return p.DisCards
}

func (p *PokerRound) GetScoreAtLeast() int {
	return p.ScoreAtLeast
}
