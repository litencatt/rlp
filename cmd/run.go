/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/litencatt/rlp/entity"
	"github.com/spf13/cobra"
)

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
	// loop
	for {
		fmt.Println("\n[Deal cards]")
		hand := deck.Deal(defaultDeal)

		// Convert to string
		var selectCards []string
		var cards []string
		for _, card := range hand {
			cards = append(cards, card.String())
		}
		promptMs := &survey.MultiSelect{
			Message: "Select cards to play",
			Options: cards,
		}
		survey.AskOne(promptMs, &selectCards, survey.WithPageSize(8))

		// Convert to Trump entity
		var selectTrumps []entity.Trump
		for _, card := range selectCards {
			// extract rank and suit from card string
			rank := strings.Split(card, " of ")[0]
			suit := strings.Split(card, " of ")[1]
			trump := entity.Trump{Rank: entity.Rank(rank), Suit: entity.Suit(suit)}
			selectTrumps = append(selectTrumps, trump)
		}

		var playOrDsicard string
		prompt := &survey.Select{
			Message: "Play or Discard:",
			Options: []string{"Play", "Discard"},
		}
		survey.AskOne(prompt, &playOrDsicard)
		if playOrDsicard == "Play" {
			poker := entity.EvaluateHand(selectTrumps)
			fmt.Printf("\nHand: %s\n\n", poker)
		} else {
			// Select cards to discard
		}

		// play again?
		var playAgain string
		promptAgain := &survey.Select{
			Message: "Play again:",
			Options: []string{"Play", "Quit"},
		}
		survey.AskOne(promptAgain, &playAgain)
		if playAgain == "Play" {

		} else {
			break
		}
	}

	return nil
}

func init() {
	rootCmd.AddCommand(runCmd)
}
