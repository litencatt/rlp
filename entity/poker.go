package entity

import (
	"sort"
)

type HandType string

const (
	HighCard      HandType = "High Card"
	OnePair                = "One Pair"
	TwoPair                = "Two Pair"
	ThreeOfAKind           = "Three of a Kind"
	Straight               = "Straight"
	Flush                  = "Flush"
	FullHouse              = "Full House"
	FourOfAKind            = "Four of a Kind"
	StraightFlush          = "Straight Flush"
	RoyalFlush             = "Royal Flush"
)

// isFlush checks if all cards in the hand have the same suit.
func isFlush(hand []Trump) bool {
	if len(hand) < 5 {
		return false
	}

	firstSuit := hand[0].Suit
	for _, card := range hand[1:] {
		if card.Suit != firstSuit {
			return false
		}
	}
	return true
}

// isStraight checks if the hand forms a sequential ranking.
func isStraight(hand []Trump) bool {
	if len(hand) < 5 {
		return false
	}

	// Sort the hand by rank (implement the sorting logic based on your Rank definition)
	sort.Slice(hand, func(i, j int) bool {
		// Implement the comparison logic for your Rank type
		return hand[i].Rank < hand[j].Rank
	})

	// Check for Ace-low straight
	if hand[0].Rank == Two && hand[1].Rank == Three && hand[2].Rank == Four && hand[3].Rank == Five && hand[4].Rank == Ace {
		return true
	}

	// Check for standard straight
	for i := 0; i < len(hand)-1; i++ {
		// convert Rank string to int
		rank1 := hand[i].GetRankNumber()
		rank2 := hand[i+1].GetRankNumber()
		if rank2-rank1 != 1 {
			return false
		}
	}
	return true
}

// groupByRank groups cards by their ranks and returns a map of rank to count.
func groupByRank(hand []Trump) map[Rank]int {
	rankCount := make(map[Rank]int)
	for _, card := range hand {
		rankCount[card.Rank]++
	}
	return rankCount
}

// EvaluateHand evaluates the given hand and returns the HandType.
func EvaluateHand(hand []Trump) HandType {
	isFlush := isFlush(hand)
	isStraight := isStraight(hand)

	if isFlush && isStraight {
		// Check for Royal Flush
		if hand[0].Rank == Ten {
			return RoyalFlush
		}
		return StraightFlush
	}

	rankCount := groupByRank(hand)
	var pairs, threes, fours int
	for _, count := range rankCount {
		switch count {
		case 2:
			pairs++
		case 3:
			threes++
		case 4:
			fours++
		}
	}

	if fours == 1 {
		return FourOfAKind
	} else if threes == 1 && pairs == 1 {
		return FullHouse
	} else if isFlush {
		return Flush
	} else if isStraight {
		return Straight
	} else if threes == 1 {
		return ThreeOfAKind
	} else if pairs == 2 {
		return TwoPair
	} else if pairs == 1 {
		return OnePair
	}

	return HighCard
}
