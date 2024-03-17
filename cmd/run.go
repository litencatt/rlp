/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/litencatt/rlp/entity"
	"github.com/manifoldco/promptui"
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
	//deck.Shuffle()

	// loop
	for {
		fmt.Println("\n[Deal 5 cards]")
		hand := deck.Deal(5)
		for _, card := range hand {
			fmt.Printf("%s of %s\n", card.Rank, card.Suit)
		}

		// Selct play or discard
		poker := entity.EvaluateHand(hand)
		fmt.Printf("\nHand: %s\n\n", poker)

		// play again?
		validate := func(input string) error {
			if input != "y" && input != "n" {
				return fmt.Errorf("invalid input(%s)", input)
			}
			return nil
		}
		p := promptui.Prompt{
			Label:    "Play again?",
			Validate: validate,
		}
		result, err := p.Run()
		if err != nil {
			return err
		}
		if result != "y" {
			break
		}
	}

	return nil
}

func init() {
	rootCmd.AddCommand(runCmd)
}
