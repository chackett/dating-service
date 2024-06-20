package rankingservice

import (
	"github.com/chackett/dating-service/repository"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewRankedResultSet_AddMatches(t *testing.T) {

	// A very rudimentary test to sanity check my "insertion sort" approach to ranking results

	rrs := NewRankedResultSet()

	rms := []RankedMatch{
		{
			User: repository.User{
				ID: 2,
			},
			Ranking: 5,
		},
		{
			User: repository.User{
				ID: 3,
			},
			Ranking: 1,
		},
		{
			User: repository.User{
				ID: 4,
			},
			Ranking: 8,
		},
	}

	for _, rm := range rms {
		rrs.AddMatch(rm)
	}

	// Just manually create expected slice in correct order
	expected := []RankedMatch{rms[2], rms[0], rms[1]}
	assert.Equal(t, expected, rrs.Matches)
}
