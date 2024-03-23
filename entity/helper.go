package entity

import "sort"

func Contains(trumps []Trump, trump Trump) bool {
	for _, t := range trumps {
		if t == trump {
			return true
		}
	}
	return false
}

func SortHand(hand []Trump) {
	sort.Slice(hand, func(i, j int) bool {
		return hand[i].GetSortOrder() < hand[j].GetSortOrder()
	})
}
