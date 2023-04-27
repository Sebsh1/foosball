package rating

import (
	"foosball/internal/user"
	"math"
)

var (
	kFactor     = 32.0
	scaleFactor = 400.0
)

func (s *ServiceImpl) calculateNewRatingsElo(winners, losers []*user.User) ([]int, []int) {
	ratingA := s.getAverageRating(winners)
	ratingB := s.getAverageRating(losers)

	return s.getNewRatings(ratingA, ratingB, winners, losers)
}

func (s *ServiceImpl) calculateNewRatingsWeighted(winners, losers []*user.User) ([]int, []int) {
	panic("unimplemented") // TODO
}

func (s *ServiceImpl) calculateNewRatingsRMS(winners, losers []*user.User) ([]int, []int) {
	ratingA := s.getRMSRating(winners)
	ratingB := s.getRMSRating(losers)

	return s.getNewRatings(ratingA, ratingB, winners, losers)
}

func (s *ServiceImpl) calculateNewRatingsGlicko2(winners, losers []*user.User) ([]int, []int) {
	panic("unimplemented") // TODO
}

func (s *ServiceImpl) getAverageRating(users []*user.User) float64 {
	sum := 0.0
	for _, u := range users {
		sum += float64(u.Rating)
	}

	return sum / float64(len(users))
}

func (s *ServiceImpl) getRMSRating(users []*user.User) float64 {
	n := 15.0
	sum := 0.0
	for _, u := range users {
		sum += math.Pow(float64(u.Rating), n)
	}
	rating := math.Pow(sum, 1/n) / float64(len(users))

	return rating
}

func (s *ServiceImpl) getNewRatings(winnersRating, losersRating float64, winners, losers []*user.User) ([]int, []int) {
	newRatingWinners := make([]int, len(winners))
	probabilityWinA := 1 / (1 + math.Pow(10, (winnersRating-losersRating)/scaleFactor))
	for i, u := range winners {
		newRatingWinners[i] = u.Rating + int(kFactor*(1-probabilityWinA))
	}

	newRatingLosers := make([]int, len(losers))
	probabilityWinB := 1 / (1 + math.Pow(10, (losersRating-winnersRating)/scaleFactor))
	for i, u := range losers {
		newRatingLosers[i] = u.Rating + int(kFactor*(-probabilityWinB))
	}

	return newRatingWinners, newRatingLosers
}
