package rlp

import (
	"strings"

	"github.com/litencatt/rlp/domain"
	"github.com/litencatt/rlp/entity"
)

type PokerRound struct {
	Deck       entity.Deck
	TotalScore int
	Hands      []entity.Trump
	Remain     []entity.Trump
	Hand       []entity.Trump
}

func NewRound() *PokerRound {
	deck := entity.NewDeck()
	deck.Shuffle()

	return &PokerRound{
		Deck: deck,
	}
}

func (p *PokerRound) Start(defaultDeal int) []string {
	p.Hands = p.Deck.Draw(defaultDeal)

	return p.HandCardString()
}

func (p *PokerRound) DrawCard(drawNum int) {
	p.Hands = nil

	drawCards := p.Deck.Draw(drawNum)
	p.Hands = append(p.Hands, p.Remain...)
	p.Hands = append(p.Hands, drawCards...)

	domain.SortHand(p.Hands)
}

func (p *PokerRound) HandCardString() []string {
	var cards []string
	for _, card := range p.Hands {
		cards = append(cards, card.String())
	}
	return cards
}

func (p *PokerRound) RemainCardString() []string {
	var cards []string
	for _, card := range p.Remain {
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
		for _, c := range p.Hands {
			if string(c.Rank) == rank && string(c.Suit) == suit {
				selectCards = append(selectCards, c)
				break
			}
		}
	}
	p.Hand = selectCards

	// Calc the remain cards
	var remainCards []entity.Trump
	for _, card := range p.Hands {
		if !domain.Contains(selectCards, card) {
			remainCards = append(remainCards, card)
		}
	}
	p.Remain = remainCards

	return len(p.Hand)
}

func (p *PokerRound) PlayHand() (entity.HandType, int) {
	handType := entity.EvaluateHand(p.Hand)
	score := entity.GetScore(handType)
	p.TotalScore += score

	return handType, score
}

func (p *PokerRound) GetTotalScore() int {
	return p.TotalScore
}
