/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/litencatt/rlp/entity"
	"github.com/spf13/cobra"
)

var debugMode bool

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run rogue-like poker",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := run(); err != nil {
			return err
		}
		return nil
	},
}

func run() error {
	fmt.Println("[Game start]")
	deck := entity.NewDeck()
	deck.Shuffle()

	defaultDeal := 8
	selectCardNum := 0
	nextDrawNum := defaultDeal
	var remainCards []entity.Trump
	totalScore := 0

	// loop
	for {
		if selectCardNum != 0 {
			nextDrawNum = selectCardNum
		}

		fmt.Println("[Draw", nextDrawNum, "cards]")
		drawCards := deck.Draw(nextDrawNum)

		var hand []entity.Trump
		if debugMode {
			fmt.Println(remainCards)
			fmt.Println(drawCards)
		}
		hand = append(hand, remainCards...)
		hand = append(hand, drawCards...)
		// sort hand
		entity.SortHand(hand)

		// Convert to string
		var selectCards []string
		var cards []string
		for _, card := range hand {
			cards = append(cards, card.String())
		}
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
			fmt.Println("Please select less than 5 cards\n")
		}

		// Convert to Trump entity
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
		// Find the remain cards
		remainCards = nil
		for _, card := range hand {
			if !entity.Contains(selectTrumps, card) {
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
			pokerHand := entity.EvaluateHand(selectTrumps)
			score := entity.GetHandScore(pokerHand)
			totalScore += score
			fmt.Printf("\nHand: %s, Score: %d\n\n", pokerHand, score)
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

func init() {
	rootCmd.AddCommand(runCmd)

	// debug mode flag
	runCmd.Flags().BoolVarP(&debugMode, "debug", "d", false, "show detail logs")
}
