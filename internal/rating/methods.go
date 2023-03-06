package rating

import (
	"foosball/internal/models"
	"math"
)

var (
	kFactor     = 32.0
	scaleFactor = 400.0
)

func (s *ServiceImpl) calculateRatingChangesElo(teamA, teamB *models.Team, winner Team) (newRatingsTeamA, newRatingsTeamB []int) {
	ratingA := s.getAverageRating(teamA)
	ratingB := s.getAverageRating(teamB)

	return s.getNewRatings(ratingA, ratingB, teamA, teamB, winner)
}

func (s *ServiceImpl) calculateRatingChangesRMS(teamA, teamB *models.Team, winner Team) (teamAChange, teamBChange []int) {
	ratingA := s.getRMSRating(teamA)
	ratingB := s.getRMSRating(teamB)

	return s.getNewRatings(ratingA, ratingB, teamA, teamB, winner)
}

func (s *ServiceImpl) getAverageRating(team *models.Team) float64 {
	sum := 0.0
	for _, p := range team.Players {
		sum += float64(p.Rating)
	}

	return sum / float64(len(team.Players))
}

func (s *ServiceImpl) getRMSRating(team *models.Team) float64 {
	n := 15.0
	sum := 0.0
	for _, p := range team.Players {
		sum += math.Pow(float64(p.Rating), n)
	}
	rating := math.Pow(sum, 1/n) / float64(len(team.Players))

	return rating
}

func (s *ServiceImpl) getNewRatings(teamARating, teamBRating float64, teamA, teamB *models.Team, winner Team) (newRatingsTeamA, newRatingsTeamB []int) {
	teamAWinner := 0.0
	if winner == TeamA {
		teamAWinner = 1.0
	}

	probabilityWinA := 1 / (1 + math.Pow(10, (teamARating-teamBRating)/scaleFactor))
	probabilityWinB := 1 / (1 + math.Pow(10, (teamBRating-teamARating)/scaleFactor))

	newRatingA := make([]int, len(teamA.Players))
	for i, p := range teamA.Players {
		newRatingA[i] = p.Rating + int(kFactor*(teamAWinner-probabilityWinA))
	}

	newRatingB := make([]int, len(teamB.Players))
	for i, p := range teamB.Players {
		newRatingB[i] = p.Rating + int(kFactor*(1-teamAWinner-probabilityWinB))
	}

	return newRatingA, newRatingB
}
