package rating

import (
	"foosball/internal/team"
	"math"
)

var (
	kFactor     = 32.0
	scaleFactor = 400.0
)

func (s *ServiceImpl) calculateRatingChangesElo(winners, losers *team.Team) ([]int, []int) {
	ratingA := s.getAverageRating(winners)
	ratingB := s.getAverageRating(losers)

	return s.getNewRatings(ratingA, ratingB, winners, losers)
}

func (s *ServiceImpl) calculateRatingChangesRMS(winners, losers *team.Team) ([]int, []int) {
	ratingA := s.getRMSRating(winners)
	ratingB := s.getRMSRating(losers)

	return s.getNewRatings(ratingA, ratingB, winners, losers)
}

func (s *ServiceImpl) getAverageRating(team *team.Team) float64 {
	sum := 0.0
	for _, p := range team.Players {
		sum += float64(p.Rating)
	}

	return sum / float64(len(team.Players))
}

func (s *ServiceImpl) getRMSRating(team *team.Team) float64 {
	n := 15.0
	sum := 0.0
	for _, p := range team.Players {
		sum += math.Pow(float64(p.Rating), n)
	}
	rating := math.Pow(sum, 1/n) / float64(len(team.Players))

	return rating
}

func (s *ServiceImpl) getNewRatings(winnersRating, losersRating float64, winners, losers *team.Team) ([]int, []int) {
	newRatingWinners := make([]int, len(winners.Players))
	probabilityWinA := 1 / (1 + math.Pow(10, (winnersRating-losersRating)/scaleFactor))
	for i, p := range winners.Players {
		newRatingWinners[i] = p.Rating + int(kFactor*(1-probabilityWinA))
	}

	newRatingLosers := make([]int, len(losers.Players))
	probabilityWinB := 1 / (1 + math.Pow(10, (losersRating-winnersRating)/scaleFactor))
	for i, p := range losers.Players {
		newRatingLosers[i] = p.Rating + int(kFactor*(-probabilityWinB))
	}

	return newRatingWinners, newRatingLosers
}
