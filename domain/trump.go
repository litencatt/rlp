package domain

import (
	"sort"

	"github.com/litencatt/rlp/entity"
)

func Contains(trumps []entity.Trump, trump entity.Trump) bool {
	for _, t := range trumps {
		if t == trump {
			return true
		}
	}
	return false
}

func SortHand(hand []entity.Trump) {
	sort.Slice(hand, func(i, j int) bool {
		return hand[i].GetSortOrder() < hand[j].GetSortOrder()
	})
}
