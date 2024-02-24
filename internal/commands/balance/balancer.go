package balance

import (
	"cmp"
	"slices"
)

func CreateTeamsGreedy(players []*Player) (*Team, *Team) {
	t1 := &Team{}
	t2 := &Team{}

	var t1Rating uint = 0
	var t2Rating uint = 0

	slices.SortFunc(players, func(a, b *Player) int {
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

func CreateTeamsBruteForce(players []*Player) (*Team, *Team) {
	matches := [][]*Team{}

	permutations := permute(players)

	for _, players := range permutations {
		midpoint := len(players) / 2

		firstHalf := players[:midpoint]
		secondHalf := players[midpoint:]

		t1 := NewTeam(firstHalf)
		t2 := NewTeam(secondHalf)

		matches = append(matches, []*Team{t1, t2})
	}

	var currentMax *int = nil
	var bestMatch []*Team = nil

	for _, teams := range matches {
		team1 := teams[0]
		team2 := teams[1]

		diff := abs(int(team1.ELO) - int(team2.ELO))

		if currentMax == nil {
			currentMax = &diff
			bestMatch = []*Team{team1, team2}
			continue
		}

		if diff < *currentMax {
			currentMax = &diff
			bestMatch = []*Team{team1, team2}
		}
	}

	return bestMatch[0], bestMatch[1]
}

// permute generates all permutations of the elements in a slice of structs
func permute(data []*Player) [][]*Player {
	var result [][]*Player
	permuteHelper(data, 0, &result)
	return result
}

// permuteHelper generates permutations recursively
func permuteHelper(data []*Player, start int, result *[][]*Player) {
	if start == len(data)-1 {
		// Make a copy of the current permutation and append it to the result
		perm := make([]*Player, len(data))
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
