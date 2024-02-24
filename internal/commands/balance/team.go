package balance

func NewTeam(players []*Player) *Team {
	var elo uint = 0

	for _, player := range players {
		elo += player.Rating
	}

	return &Team{
		Players: players,
		ELO:     elo,
	}
}

type Team struct {
	Players []*Player
	ELO     uint
}
