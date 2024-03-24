package rlp

import (
	"fmt"
	"os"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
)

const Name string = "rlp"

var Version = "dev"

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
	fmt.Println("Welcome to Rogue-like Poker!")
	fmt.Println("Round start")
	fmt.Println()

	round := NewRound()

	var selectCardNum int
	nextDrawNum := r.DefaultDeal
	for {
		if selectCardNum != 0 {
			nextDrawNum = selectCardNum
		}

		if r.DebugMode {
			fmt.Println("[Draw", nextDrawNum, "cards]")
			fmt.Println()
		}

		round.DrawCard(nextDrawNum)
		handCards := round.HandCardString()

		fmt.Printf("Score at least: %d\n", round.ScoreAtLeast)
		fmt.Printf("Round score: %d\n", round.TotalScore)
		fmt.Printf("Hands: %d, Discards: %d\n", round.Hands, round.Discards)
		fmt.Println()

		time.Sleep(1 * time.Second)

		// Select cards
		var selectCards []string
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
		var pdOptions = []string{"Play"}
		if round.Discards > 0 {
			pdOptions = append(pdOptions, "Discard")
		}
		prompt := &survey.Select{
			Message: "Play or Discard:",
			Options: pdOptions,
		}
		if err := survey.AskOne(prompt, &playOrDsicard); err == terminal.InterruptErr {
			fmt.Println("interrupted")
			os.Exit(0)
		}
		if playOrDsicard == "Play" {
			round.Hands--
			handType, score := round.PlayHand()
			fmt.Printf("\nHand: %s, Score: %d\n\n", handType, score)

			time.Sleep(1 * time.Second)
		} else {
			round.Discards--
		}

		// show remain cards
		if r.DebugMode {
			fmt.Print("[Remain cards]\n")
			remainCards := round.RemainCardString()
			for _, card := range remainCards {
				fmt.Println(card)
			}
			fmt.Println()
		}

		if playOrDsicard == "Discard" {
			continue
		}

		if round.TotalScore >= round.ScoreAtLeast {
			fmt.Printf("Score at least: %d, Round score: %d\n", round.ScoreAtLeast, round.TotalScore)
			fmt.Println("You win!")
			break
		} else if round.Hands <= 0 {
			fmt.Printf("Score at least: %d, Round score: %d\n", round.ScoreAtLeast, round.TotalScore)
			fmt.Println("You lose!")
			break
		} else {
			if round.Hands > 0 && round.TotalScore < round.ScoreAtLeast {
				continue
			}
		}
	}

	return nil
}
