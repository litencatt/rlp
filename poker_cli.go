package pkr

import (
	"fmt"
	"os"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/litencatt/pkr/entity"
	"github.com/litencatt/pkr/service"
)

const Name string = "pkr"

var Version = "dev"

type PokerCLI struct {
	DebugMode bool
	service   service.PokerService
}

func NewPokerCLI() *PokerCLI {
	return &PokerCLI{
		service: service.NewPokerService(service.PokerServiceConfig{
			DebugMode: false,
		}),
	}
}

func (cli *PokerCLI) Run() error {
	fmt.Println("Welcome to Poker!")
	fmt.Println()

BlindLoop:
	for _, ante := range entity.Antes {
		for _, blind := range entity.Blinds {
			fmt.Printf("\n\nRound start\nAnte:%d, Blind:%v\n\n", ante, blind)

			ScoreAtLeast := int(float64(ante) * blind)
			fmt.Printf("Score at least: %d\n", ScoreAtLeast)
			cli.service.StartRound(ScoreAtLeast)

			var selectCardNum int
			var nextDrawNum int

			for {
				if selectCardNum == 0 {
					nextDrawNum = cli.service.GetNextDrawNum()
				} else {
					// Draw cards num is same as the last selected cards num
					nextDrawNum = selectCardNum
				}

				if err := cli.service.DrawCard(nextDrawNum); err != nil {
					return err
				}

				round := cli.service.GetRoundInfo()
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
					chip, mult := cli.service.GetChipAndMult(handType, 1)
					fmt.Printf("\nHand: %s(Chip: %d, Mult: %d)\n", handType, chip, mult)

					// get card rank and add to chip
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
				if cli.DebugMode {
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
					break BlindLoop
				} else {
					if round.Hands > 0 && round.TotalScore < round.ScoreAtLeast {
						continue
					}
				}
			}
		}
	}

	return nil
}
