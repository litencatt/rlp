package entity

import "sort"

type PokerHandCard struct {
	Trumps []Trump
}

func (p *PokerHandCard) Contains(trump Trump) bool {
	for _, t := range p.Trumps {
		if t == trump {
			return true
		}
	}
	return false
}

func (p *PokerHandCard) Sort() {
	sort.Slice(p.Trumps, func(i, j int) bool {
		return p.Trumps[i].GetSortOrder() < p.Trumps[j].GetSortOrder()
	})
}
