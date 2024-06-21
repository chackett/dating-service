package rankingservice

import (
	"github.com/chackett/dating-service/repository"
	"sort"
)

type RankedMatch struct {
	repository.User
	Ranking        int `json:"ranking"`
	DistanceFromMe int `json:"distanceFromMe"`
}

type RankedResultSet struct {
	Matches []RankedMatch `json:"matches,omitempty"`
}

func NewRankedResultSet() RankedResultSet {
	result := RankedResultSet{
		Matches: make([]RankedMatch, 0),
	}
	return result
}

func (r *RankedResultSet) AddMatch(input RankedMatch) {
	index := sort.Search(len(r.Matches), func(i int) bool {
		return r.Matches[i].Ranking <= input.Ranking
	})

	r.Matches = append(r.Matches, RankedMatch{})
	copy(r.Matches[index+1:], r.Matches[index:])
	r.Matches[index] = input
}
