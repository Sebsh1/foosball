package rating

import (
	"math"
)

type Method string

const (
	Elo     Method = "elo"
	RMS     Method = "rms"
	Glicko2 Method = "glicko2"
)

var (
	kFactor     = 32.0
	scaleFactor = 400.0
)

func (s *ServiceImpl) calculateNewRatingsElo(winners, losers []Rating) []Rating {
	ratingWinners := s.getAverageRating(winners)
	ratingLosers := s.getAverageRating(losers)

	return s.getNewRatings(ratingWinners, ratingLosers, winners, losers)
}

func (s *ServiceImpl) calculateNewRatingsRMS(winners, losers []Rating) []Rating {
	ratingWinners := s.getRMSRating(winners)
	ratingLosers := s.getRMSRating(losers)

	return s.getNewRatings(ratingWinners, ratingLosers, winners, losers)
}

func (s *ServiceImpl) calculateNewRatingsGlicko2(winners, losers []Rating) []Rating {
	panic("unimplemented") // TODO
}

func (s *ServiceImpl) getNewRatings(winnersRating, losersRating float64, winners, losers []Rating) []Rating {
	newRatingWinners := make([]Rating, len(winners))
	probabilityWinA := 1 / (1 + math.Pow(10, (winnersRating-losersRating)/scaleFactor))
	for i, r := range winners {
		r.Rating = r.Rating + int(kFactor*(1-probabilityWinA))
		newRatingWinners[i] = r
	}

	newRatingLosers := make([]Rating, len(losers))
	probabilityWinB := 1 / (1 + math.Pow(10, (losersRating-winnersRating)/scaleFactor))
	for i, r := range losers {
		r.Rating = r.Rating + int(kFactor*(-probabilityWinB))
		newRatingLosers[i] = r
	}

	newRatings := append(newRatingWinners, newRatingLosers...)
	return newRatings
}

func (s *ServiceImpl) getAverageRating(ratings []Rating) float64 {
	sum := 0
	for _, u := range ratings {
		sum += u.Rating
	}

	return float64(sum) / float64(len(ratings))
}

func (s *ServiceImpl) getRMSRating(ratings []Rating) float64 {
	n := 15.0
	sum := 0.0
	for _, u := range ratings {
		sum += math.Pow(float64(u.Rating), n)
	}

	return math.Pow(sum, 1/n) / float64(len(ratings))
}
