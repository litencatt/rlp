package rlp

import (
	"fmt"
	"os"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/litencatt/rlp/domain"
	"github.com/litencatt/rlp/entity"
)

type RogurLikePoker struct {
	Deck        entity.Deck
	DefaultDeal int
	DebugMode   bool
}

func NewRogurLikePoker() *RogurLikePoker {
	deck := entity.NewDeck()

	return &RogurLikePoker{
		Deck:        deck,
		DefaultDeal: 8,
	}
}

func (r *RogurLikePoker) Run() error {
	r.Deck.Shuffle()
	r.DefaultDeal = 8
	selectCardNum := 0
	nextDrawNum := r.DefaultDeal

	var remainCards []entity.Trump
	totalScore := 0

	// loop
	for {
		if selectCardNum != 0 {
			nextDrawNum = selectCardNum
		}

		fmt.Println("[Draw", nextDrawNum, "cards]")
		drawCards := r.Deck.Draw(nextDrawNum)

		var hand []entity.Trump
		if r.DebugMode {
			fmt.Println(remainCards)
			fmt.Println(drawCards)
		}
		hand = append(hand, remainCards...)
		hand = append(hand, drawCards...)
		// sort hand
		domain.SortHand(hand)

		// Convert to string
		var selectCards []string
		var cards []string
		for _, card := range hand {
			cards = append(cards, card.String())
		}
		// Select cards
		for {
			selectCards = nil
			promptMs := &survey.MultiSelect{
				Message: "Select cards",
				Options: cards,
			}
			err := survey.AskOne(promptMs, &selectCards, survey.WithPageSize(8))
			if err == terminal.InterruptErr {
				fmt.Println("interrupted")
				os.Exit(0)
			}
			selectCardNum = len(selectCards)
			if selectCardNum <= 5 {
				break
			}
			fmt.Println("Please select less than 5 cards")
			fmt.Println()
		}

		// Convert select cards to Trump entity
		var selectTrumps []entity.Trump
		for _, card := range selectCards {
			// extract rank and suit from card string
			rank := strings.Split(card, " of ")[0]
			suit := strings.Split(card, " of ")[1]
			// Find the card from hand
			for _, c := range hand {
				if string(c.Rank) == rank && string(c.Suit) == suit {
					selectTrumps = append(selectTrumps, c)
					break
				}
			}
		}
		// Calc the remain cards
		remainCards = nil
		for _, card := range hand {
			if !domain.Contains(selectTrumps, card) {
				remainCards = append(remainCards, card)
			}
		}

		var playOrDsicard string
		prompt := &survey.Select{
			Message: "Play or Discard:",
			Options: []string{"Play", "Discard"},
		}
		if err := survey.AskOne(prompt, &playOrDsicard); err == terminal.InterruptErr {
			fmt.Println("interrupted")
			os.Exit(0)
		}
		if playOrDsicard == "Play" {
			handType := entity.EvaluateHand(selectTrumps)
			score := entity.GetScore(handType)
			totalScore += score
			fmt.Printf("\nHand: %s, Score: %d\n\n", handType, score)
		}

		fmt.Print("Remain cards:\n")
		for _, card := range remainCards {
			fmt.Println(card)
		}
		fmt.Println()

		// play again?
		var playAgain string
		promptAgain := &survey.Select{
			Message: "Play again? (Total score: " + fmt.Sprintf("%d)", totalScore),
			Options: []string{"Play", "Quit"},
		}
		if err := survey.AskOne(promptAgain, &playAgain); err == terminal.InterruptErr {
			fmt.Println("interrupted")
			os.Exit(0)
		}
		if playAgain == "Play" {
			continue
		} else {
			break
		}
	}

	return nil
}
