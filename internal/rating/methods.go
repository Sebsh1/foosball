package rating

import "foosball/internal/player"

func (s *ServiceImpl) calculateRatingChangesElo(teamA []*player.Player, teamB []*player.Player) (teamAChange []int, teamBChange []int) {
	// TODO
	return make([]int, len(teamA), 0), make([]int, len(teamB), 0)
}

func (s *ServiceImpl) calculateRatingChangesWeighted(teamA []*player.Player, teamB []*player.Player) (teamAChange []int, teamBChange []int) {
	// TODO
	return make([]int, len(teamA), 0), make([]int, len(teamB), 0)
}

func (s *ServiceImpl) calculateRatingChangesRMS(teamA []*player.Player, teamB []*player.Player) (teamAChange []int, teamBChange []int) {
	// TODO
	return make([]int, len(teamA), 0), make([]int, len(teamB), 0)
}
