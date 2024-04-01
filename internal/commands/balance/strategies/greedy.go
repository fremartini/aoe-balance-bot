package strategies

import (
	"aoe-bot/internal/commands/balance"
	"cmp"
	"slices"
)

type greedyStrategy struct{}

func NewGreedy() *greedyStrategy {
	return &greedyStrategy{}
}

func (*greedyStrategy) CreateTeams(players []*balance.Player) (*balance.Team, *balance.Team) {
	t1 := &balance.Team{}
	t2 := &balance.Team{}

	var t1Rating uint = 0
	var t2Rating uint = 0

	slices.SortFunc(players, func(a, b *balance.Player) int {
		return cmp.Compare(b.Rating, a.Rating)
	})

	for _, player := range players {
		if t1Rating < t2Rating {
			t1.Players = append(t1.Players, player)
			t1Rating += player.Rating
		} else {
			t2.Players = append(t2.Players, player)
			t2Rating += player.Rating
		}
	}

	t1.ELO = t1Rating
	t2.ELO = t2Rating

	return t1, t2
}
