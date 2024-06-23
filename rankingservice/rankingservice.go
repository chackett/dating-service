package rankingservice

import (
	"github.com/chackett/dating-service/repository"
	"sort"
)

// RankedMatch is a type composes of User, but extended with fields specific to the context of a "match"
// Such as the profiles match against user and distance from the user.
type RankedMatch struct {
	repository.User
	// Ranking is the score of how well matched the profile is to the user.
	Ranking int `json:"ranking"`
	// DistanceFromMe specifies distance in KM from the user
	DistanceFromMe int `json:"distanceFromMe"`
}

// RankedResultSet set of results to be returned to user
type RankedResultSet struct {
	Matches []RankedMatch `json:"matches,omitempty"`
}

func NewRankedResultSet() RankedResultSet {
	result := RankedResultSet{
		Matches: make([]RankedMatch, 0),
	}
	return result
}

// AddMatch is used to insert matches into result set in a sorted fashion, based on the ranking in the profile.
// This has the effect of returning sorted results to the user.
func (r *RankedResultSet) AddMatch(input RankedMatch) {
	index := sort.Search(len(r.Matches), func(i int) bool {
		return r.Matches[i].Ranking <= input.Ranking
	})

	r.Matches = append(r.Matches, RankedMatch{})
	copy(r.Matches[index+1:], r.Matches[index:])
	r.Matches[index] = input
}
