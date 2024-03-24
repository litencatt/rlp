package rlp

import (
	"fmt"
	"os"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/litencatt/rlp/entity"
)

const Name string = "rlp"

var Version = "dev"

type RogurLikePoker struct {
	DebugMode bool
	RunInfo   *entity.RunInfo
}

func NewRogurLikePoker() *RogurLikePoker {
	return &RogurLikePoker{
		RunInfo: entity.NewRunInfo(),
	}
}

func (r *RogurLikePoker) Run() error {
	fmt.Println("Welcome to Rogue-like Poker!")
	fmt.Println("Round start")
	fmt.Println()

	// Select small blind
	ante := r.RunInfo.Ante
	blind := ante.Blinds[0]
	ScoreAtLeast := int(float64(ante.GetAnteBase()) * blind.Multi)

	round := entity.PokerRound{
		Deck:         r.RunInfo.Deck,
		TotalScore:   0,
		Hands:        r.RunInfo.Hands,
		Discards:     r.RunInfo.Discards,
		ScoreAtLeast: ScoreAtLeast,
	}

	round.Deck.Shuffle()

	var selectCardNum int
	var nextDrawNum int
	for {
		if selectCardNum == 0 {
			nextDrawNum = r.RunInfo.DefaultDeal
		} else {
			// Draw cards num is same as the last selected cards num
			nextDrawNum = selectCardNum
		}

		drawCards := round.DrawCard(nextDrawNum)
		if r.DebugMode {
			fmt.Println("[Draw", nextDrawNum, "cards]")
			for _, card := range drawCards {
				fmt.Println(card.String())
			}
			fmt.Println()
		}

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
				Options: round.HandCardString(),
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

		round.SetSelectCards(selectCards)
		fmt.Print("[Selected cards]\n")
		for _, card := range selectCards {
			fmt.Println(card)
		}
		fmt.Println()

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

			handType := round.PlayHand()
			chip, mult := r.RunInfo.PokerHands.GetChipAndMult(handType, 1)
			fmt.Printf("\nHand: %s(Chip: %d, Mult: %d)\n", handType, chip, mult)

			// get card rank
			handsRankTotal := round.GetSelectCardsRankTotal()
			fmt.Printf("Hands rank total: %d\n", handsRankTotal)
			chip += handsRankTotal

			score := chip * mult
			fmt.Printf("\nChip: %d, Mult: %d\n", chip, mult)
			fmt.Printf("\nScore: %d\n\n", score)

			round.TotalScore += score

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
