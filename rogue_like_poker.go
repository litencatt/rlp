package rlp

import (
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
)

type RogurLikePoker struct {
	DefaultDeal int
	DebugMode   bool
}

func NewRogurLikePoker() *RogurLikePoker {
	return &RogurLikePoker{
		DefaultDeal: 8,
	}
}

func (r *RogurLikePoker) Run() error {
	round := NewRound()

	var selectCardNum int
	nextDrawNum := r.DefaultDeal
	for {
		if selectCardNum != 0 {
			nextDrawNum = selectCardNum
		}

		fmt.Println("[Draw", nextDrawNum, "cards]")
		round.DrawCard(nextDrawNum)
		handCards := round.HandCardString()

		var selectCards []string
		// Select cards
		for {
			selectCards = nil
			promptMs := &survey.MultiSelect{
				Message: "Select cards",
				Options: handCards,
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
		round.SelectCards(selectCards)

		// Play or Discard
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
			handType, score := round.PlayHand()
			fmt.Printf("\nHand: %s, Score: %d\n\n", handType, score)
		}

		// show remain cards
		fmt.Print("Remain cards:\n")
		remainCards := round.RemainCardString()
		for _, card := range remainCards {
			fmt.Println(card)
		}
		fmt.Println()

		if playOrDsicard == "Discard" {
			continue
		}

		// play again?
		var playAgain string
		promptAgain := &survey.Select{
			Message: "Play again? (Total score: " + fmt.Sprintf("%d)", round.GetTotalScore()),
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
