package rating

import (
	"foosball/internal/models"
	"math"
)

var (
	kFactor     = 32.0
	scaleFactor = 400.0
)

func (s *ServiceImpl) calculateRatingChangesElo(teamA, teamB []*models.Player, winner Team) (newRatingsTeamA, newRatingsTeamB []int) {
	ratingA := s.getAverageRating(teamA)
	ratingB := s.getAverageRating(teamB)

	return s.getNewRatings(ratingA, ratingB, teamA, teamB, winner)
}

func (s *ServiceImpl) calculateRatingChangesRMS(teamA, teamB []*models.Player, winner Team) (teamAChange, teamBChange []int) {
	ratingA := s.getRMSRating(teamA)
	ratingB := s.getRMSRating(teamB)

	return s.getNewRatings(ratingA, ratingB, teamA, teamB, winner)
}

func (s *ServiceImpl) getAverageRating(players []*models.Player) float64 {
	sum := 0.0
	for _, p := range players {
		sum += float64(p.Rating)
	}

	return sum / float64(len(players))
}

func (s *ServiceImpl) getRMSRating(players []*models.Player) float64 {
	n := 15.0
	sum := 0.0
	for _, p := range players {
		sum += math.Pow(float64(p.Rating), n)
	}
	rating := math.Pow(sum, 1/n) / float64(len(players))

	return rating
}

func (s *ServiceImpl) getNewRatings(teamARating, teamBRating float64, teamA, teamB []*models.Player, winner Team) (newRatingsTeamA, newRatingsTeamB []int) {
	teamAWinner := 0.0
	if winner == TeamA {
		teamAWinner = 1.0
	}

	probabilityWinA := 1 / (1 + math.Pow(10, (teamARating-teamBRating)/scaleFactor))
	probabilityWinB := 1 / (1 + math.Pow(10, (teamBRating-teamARating)/scaleFactor))

	newRatingA := make([]int, len(teamA))
	for i, p := range teamA {
		newRatingA[i] = p.Rating + int(kFactor*(teamAWinner-probabilityWinA))
	}

	newRatingB := make([]int, len(teamB))
	for i, p := range teamB {
		newRatingB[i] = p.Rating + int(kFactor*(1-teamAWinner-probabilityWinB))
	}

	return newRatingA, newRatingB
}
