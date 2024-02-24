package balance_test

import (
	"aoe-bot/internal/commands/balance"
	"reflect"
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
	t1, t2 := balance.CreateTeamsGreedy(players)

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

func TestBruteForce_FourPlayers_FindsFairMatch(t *testing.T) {
	// arrange
	player1 := &balance.Player{
		Name:   "p1",
		Rating: 2300,
	}

	player2 := &balance.Player{
		Name:   "p2",
		Rating: 2000,
	}

	player3 := &balance.Player{
		Name:   "p3",
		Rating: 1700,
	}

	player4 := &balance.Player{
		Name:   "p4",
		Rating: 1600,
	}

	players := []*balance.Player{player1, player2, player3, player4}

	// act
	t1, t2 := balance.CreateTeamsBruteForce(players)

	// assert
	t1Expected := []*balance.Player{
		player1,
		player4,
	}

	t2Expected := []*balance.Player{
		player3,
		player2,
	}

	if !reflect.DeepEqual(t1.Players, t1Expected) {
		t.Errorf("%v was not equal to %v", t1.Players, t1Expected)
	}

	if !reflect.DeepEqual(t2.Players, t2Expected) {
		t.Errorf("%v was not equal to %v", t2.Players, t2Expected)
	}
}

func TestBruteForce_FivePlayers_FindsFairMatch(t *testing.T) {
	// arrange
	player1 := &balance.Player{
		Name:   "p1",
		Rating: 2300,
	}

	player2 := &balance.Player{
		Name:   "p2",
		Rating: 2000,
	}

	player3 := &balance.Player{
		Name:   "p3",
		Rating: 1700,
	}

	player4 := &balance.Player{
		Name:   "p4",
		Rating: 1600,
	}

	player5 := &balance.Player{
		Name:   "p5",
		Rating: 1500,
	}

	players := []*balance.Player{player1, player2, player3, player4, player5}

	// act
	t1, t2 := balance.CreateTeamsBruteForce(players)

	// assert
	t1Expected := []*balance.Player{
		player1,
		player2,
	}

	t2Expected := []*balance.Player{
		player3,
		player4,
		player5,
	}

	if !reflect.DeepEqual(t1.Players, t1Expected) {
		t.Errorf("%v was not equal to %v", t1.Players, t1Expected)
	}

	if !reflect.DeepEqual(t2.Players, t2Expected) {
		t.Errorf("%v was not equal to %v", t2.Players, t2Expected)
	}
}

func TestBruteForce_EightPlayers_FindsFairMatch(t *testing.T) {
	// arrange
	player1 := &balance.Player{
		Name:   "p1",
		Rating: 2300,
	}

	player2 := &balance.Player{
		Name:   "p2",
		Rating: 2000,
	}

	player3 := &balance.Player{
		Name:   "p3",
		Rating: 1700,
	}

	player4 := &balance.Player{
		Name:   "p4",
		Rating: 1600,
	}

	player5 := &balance.Player{
		Name:   "p5",
		Rating: 1500,
	}

	player6 := &balance.Player{
		Name:   "p6",
		Rating: 1500,
	}

	player7 := &balance.Player{
		Name:   "p7",
		Rating: 800,
	}

	player8 := &balance.Player{
		Name:   "p8",
		Rating: 800,
	}

	players := []*balance.Player{player1, player2, player3, player4, player5, player6, player7, player8}

	// act
	t1, t2 := balance.CreateTeamsBruteForce(players)

	// assert
	t1Expected := []*balance.Player{
		player1,
		player5,
		player6,
		player7,
	}

	t2Expected := []*balance.Player{
		player2,
		player3,
		player4,
		player8,
	}

	if !reflect.DeepEqual(t1.Players, t1Expected) {
		t.Errorf("%v was not equal to %v", t1.Players, t1Expected)
	}

	if !reflect.DeepEqual(t2.Players, t2Expected) {
		t.Errorf("%v was not equal to %v", t2.Players, t2Expected)
	}
}
