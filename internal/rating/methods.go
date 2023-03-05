package rating

import (
	"foosball/internal/team"
	"math"
)

var (
	kFactor     = 32.0
	scaleFactor = 400.0
)

func (s *ServiceImpl) calculateRatingChangesElo(teamA, teamB *team.Team, winner Team) (newRatingsTeamA, newRatingsTeamB []int) {
	ratingA := s.getAverageRating(teamA)
	ratingB := s.getAverageRating(teamB)

	return s.getNewRatings(ratingA, ratingB, teamA, teamB, winner)
}

func (s *ServiceImpl) calculateRatingChangesRMS(teamA, teamB *team.Team, winner Team) (teamAChange, teamBChange []int) {
	ratingA := s.getRMSRating(teamA)
	ratingB := s.getRMSRating(teamB)

	return s.getNewRatings(ratingA, ratingB, teamA, teamB, winner)
}

func (s *ServiceImpl) getAverageRating(team *team.Team) float64 {
	sum := 0.0
	for _, p := range team.GetPlayers() {
		sum += float64(p.Rating)
	}

	return sum / float64(len(team.GetPlayers()))
}

func (s *ServiceImpl) getRMSRating(team *team.Team) float64 {
	n := 15.0
	sum := 0.0
	for _, p := range team.GetPlayers() {
		sum += math.Pow(float64(p.Rating), n)
	}
	rating := math.Pow(sum, 1/n) / float64(len(team.GetPlayers()))

	return rating
}

func (s *ServiceImpl) getNewRatings(teamARating, teamBRating float64, teamA, teamB *team.Team, winner Team) (newRatingsTeamA, newRatingsTeamB []int) {
	teamAWinner := 0.0
	if winner == TeamA {
		teamAWinner = 1.0
	}

	probabilityWinA := 1 / (1 + math.Pow(10, (teamARating-teamBRating)/scaleFactor))
	probabilityWinB := 1 / (1 + math.Pow(10, (teamBRating-teamARating)/scaleFactor))

	playersTeamA := teamA.GetPlayers()
	newRatingA := make([]int, len(playersTeamA))
	for i, p := range playersTeamA {
		newRatingA[i] = p.Rating + int(kFactor*(teamAWinner-probabilityWinA))
	}

	playersTeamB := teamB.GetPlayers()
	newRatingB := make([]int, len(playersTeamB))
	for i, p := range playersTeamB {
		newRatingB[i] = p.Rating + int(kFactor*(1-teamAWinner-probabilityWinB))
	}

	return newRatingA, newRatingB
}
