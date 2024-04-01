package strategies

import "aoe-bot/internal/commands/balance"

type bruteForceStrategy struct{}

func NewBruteForce() *bruteForceStrategy {
	return &bruteForceStrategy{}
}

func (*bruteForceStrategy) CreateTeams(players []*balance.Player) (*balance.Team, *balance.Team) {
	matches := [][]*balance.Team{}

	permutations := permute(players)

	for _, players := range permutations {
		midpoint := len(players) / 2

		firstHalf := players[:midpoint]
		secondHalf := players[midpoint:]

		t1 := balance.NewTeam(firstHalf)
		t2 := balance.NewTeam(secondHalf)

		matches = append(matches, []*balance.Team{t1, t2})
	}

	var currentMax *int = nil
	var bestMatch []*balance.Team = nil

	for _, teams := range matches {
		team1 := teams[0]
		team2 := teams[1]

		diff := abs(int(team1.ELO) - int(team2.ELO))

		if currentMax == nil {
			currentMax = &diff
			bestMatch = []*balance.Team{team1, team2}
			continue
		}

		if diff < *currentMax {
			currentMax = &diff
			bestMatch = []*balance.Team{team1, team2}
		}
	}

	return bestMatch[0], bestMatch[1]
}

// permute generates all permutations of the elements in a slice of structs
func permute(data []*balance.Player) [][]*balance.Player {
	var result [][]*balance.Player
	permuteHelper(data, 0, &result)
	return result
}

// permuteHelper generates permutations recursively
func permuteHelper(data []*balance.Player, start int, result *[][]*balance.Player) {
	if start == len(data)-1 {
		// Make a copy of the current permutation and append it to the result
		perm := make([]*balance.Player, len(data))
		copy(perm, data)
		*result = append(*result, perm)
		return
	}

	for i := start; i < len(data); i++ {
		// Swap the current element with element at index 'start'
		data[start], data[i] = data[i], data[start]

		// Recursively generate permutations for the rest of the elements
		permuteHelper(data, start+1, result)

		// Backtrack by swapping the elements back to their original positions
		data[start], data[i] = data[i], data[start]
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
