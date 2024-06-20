package rankingservice

import "github.com/chackett/dating-service/repository"

type RankedMatch struct {
	repository.User
	Ranking int `json:"ranking"`
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
	// Base case - start off empty
	if len(r.Matches) == 0 {
		r.Matches = append(r.Matches, input)
		return
	}

	for i := 0; i < len(r.Matches); i++ {
		m := r.Matches[i]
		if input.Ranking < m.Ranking || input.Ranking == m.Ranking {
			r.Matches = append(r.Matches, input)
			break
		} else if input.Ranking > m.Ranking {
			// Prepend
			r.Matches = append([]RankedMatch{input}, r.Matches...)
			break
		}
	}
}
