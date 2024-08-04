package service

import (
	"fmt"

	"github.com/litencatt/pkr/entity"
)

type PokerService interface {
	StartRound(int) error
	DrawCard(int) error
	PlayHand() error
	DiscardHand() error
	NextRound() error
	NextAnte() error
	GetNextDrawNum() int
	GetRoundInfo() *entity.PokerRound
	GetChipAndMult(entity.HandType, int) (int, int)
}

type pokerService struct {
	config  PokerServiceConfig
	runInfo *entity.RunInfo
}

func NewPokerService(config PokerServiceConfig) PokerService {
	return &pokerService{
		config:  config,
		runInfo: entity.NewRunInfo(),
	}
}

type PokerServiceConfig struct {
	DebugMode bool
}

func (s *pokerService) GetNextDrawNum() int {
	return s.runInfo.DefaultDeal
}

func (s *pokerService) GetRoundInfo() *entity.PokerRound {
	return s.runInfo.Round
}

func (s *pokerService) GetChipAndMult(handType entity.HandType, level int) (int, int) {
	return s.runInfo.PokerHands.GetChipAndMult(handType, level)
}

func (s *pokerService) StartRound(ScoreAtLeast int) error {
	scoreAtLeast := int(float64(s.runInfo.Ante) * s.runInfo.Blind)
	s.runInfo.Round = &entity.PokerRound{
		Deck:         s.runInfo.Deck,
		TotalScore:   0,
		Hands:        s.runInfo.Hands,
		Discards:     s.runInfo.Discards,
		ScoreAtLeast: scoreAtLeast,
	}

	s.runInfo.Round.Deck.Shuffle()

	return nil
}

func (s *pokerService) DrawCard(nextDrawNum int) error {
	drawCards := s.runInfo.Round.DrawCard(nextDrawNum)
	if s.config.DebugMode {
		fmt.Println("[Draw", nextDrawNum, "cards]")
		for _, card := range drawCards {
			fmt.Println(card.String())
		}
		fmt.Println()
	}

	return nil
}

func (s *pokerService) PlayHand() error {
	return nil
}

func (s *pokerService) DiscardHand() error {
	return nil
}

func (s *pokerService) NextRound() error {
	return nil
}

func (s *pokerService) NextAnte() error {
	return nil
}

// NewPokerServiceConfig returns a new PokerServiceConfig
func NewPokerServiceConfig() PokerServiceConfig {
	return PokerServiceConfig{}
}
