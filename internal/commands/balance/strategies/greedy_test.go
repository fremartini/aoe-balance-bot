package strategies_test

import (
	"aoe-bot/internal/commands/balance"
	"aoe-bot/internal/commands/balance/strategies"
	"testing"
)

func TestCreateTeamsGreedy_TwoPlayers_CreatesTeams(t *testing.T) {
	// arrange
	player1 := &balance.Player{
		Name:   "p1",
		Rating: 1000,
	}

	player2 := &balance.Player{
		Name:   "p2",
		Rating: 2000,
	}

	players := []*balance.Player{player1, player2}

	// act
	t1, t2 := strategies.NewGreedy().CreateTeams(players)

	// assert
	if t1.ELO != player1.Rating {
		t.Errorf("Incorrect team 1 ELO")
	}

	if t2.ELO != player2.Rating {
		t.Errorf("Incorrect team 2 ELO")
	}

	if t1.Players[0].Name != player1.Name {
		t.Errorf("Incorrect team 1 player name")
	}

	if t1.Players[0].Rating != player1.Rating {
		t.Errorf("Incorrect team 1 player rating")
	}

	if t2.Players[0].Name != player2.Name {
		t.Errorf("Incorrect team 1 player name")
	}

	if t2.Players[0].Rating != player2.Rating {
		t.Errorf("Incorrect team 2 player rating")
	}
}
