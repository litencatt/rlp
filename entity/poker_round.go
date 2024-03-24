package entity

import (
	"strings"
)

type PokerRound struct {
	Deck          Deck
	TotalScore    int
	HandCards     PokerHandCard
	RemainCards   []Trump
	SelectedCards []Trump
	Hands         int
	Discards      int
	ScoreAtLeast  int
}

func (p *PokerRound) DrawCard(drawNum int) []Trump {
	p.HandCards.Trumps = nil

	drawCards := p.Deck.Draw(drawNum)
	p.HandCards.Trumps = append(p.HandCards.Trumps, p.RemainCards...)
	p.HandCards.Trumps = append(p.HandCards.Trumps, drawCards...)

	p.HandCards.Sort()

	return drawCards
}

func (p *PokerRound) HandCardString() []string {
	var cards []string
	for _, card := range p.HandCards.Trumps {
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
	var selectCards []Trump
	for _, card := range cards {
		// extract rank and suit from card string
		rank := strings.Split(card, " of ")[0]
		suit := strings.Split(card, " of ")[1]
		// Find the card from hand
		for _, t := range p.HandCards.Trumps {
			if string(t.Rank) == rank && string(t.Suit) == suit {
				selectCards = append(selectCards, t)
				break
			}
		}
	}
	p.SelectedCards = selectCards

	// Calc the RemainCards cards
	var RemainCardsCards []Trump
	for _, card := range p.HandCards.Trumps {
		if !Contains(selectCards, card) {
			RemainCardsCards = append(RemainCardsCards, card)
		}
	}
	p.RemainCards = RemainCardsCards
}

func (p *PokerRound) PlayHand() (HandType, int) {
	handType := EvaluateHand(p.SelectedCards)
	score := GetScore(handType)
	p.TotalScore += score

	return handType, score
}
