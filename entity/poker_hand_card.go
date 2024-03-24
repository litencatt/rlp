package entity

import "sort"

type PokerHandCard struct {
	Trumps []Trump
}

func (p *PokerHandCard) Sort() {
	sort.Slice(p.Trumps, func(i, j int) bool {
		return p.Trumps[i].GetSortOrder() < p.Trumps[j].GetSortOrder()
	})
}
