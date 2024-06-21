package repository

import (
	"github.com/umahmood/haversine"
	"strconv"
	"strings"
	"time"
)

type User struct {
	ID          int        `json:"id,omitempty"`
	Email       string     `json:"email,omitempty"`
	Password    string     `json:"password,omitempty" `
	Name        string     `json:"name,omitempty"`
	Gender      string     `json:"gender,omitempty"`
	DateOfBirth *time.Time `json:"dateOfBirth,omitempty"`
	Age         int        `json:"age,omitempty" gorm:"-"`
	Location    string     `json:"location,omitempty"`
}

func (u *User) CalculateAge() int {
	now := time.Now()
	age := now.Year() - u.DateOfBirth.Year()

	if now.YearDay() < u.DateOfBirth.YearDay() {
		age = age - 1
	}
	return age
}

func (u *User) ReadLocation() haversine.Coord {
	spl := strings.Split(u.Location, ",")

	fLat, err := strconv.ParseFloat(spl[0], 32)
	if err != nil {
		return haversine.Coord{}
	}
	fLong, err := strconv.ParseFloat(spl[1], 32)
	if err != nil {
		return haversine.Coord{}
	}
	return haversine.Coord{
		Lat: fLat,
		Lon: fLong,
	}
}

func (u *User) DistanceFromUser(candidate User) int {
	_, km := haversine.Distance(u.ReadLocation(), candidate.ReadLocation())

	return int(km)
}

func (u *User) RankCandidate(candidate User, userPrefs UserPreferences, canPrefs UserPreferences) (int, error) {

	// Ranking is used to score a candidate. Note for total mismatched candidates, -1 is returned immediately.
	// While other comparisons might not be a direct match, it doesn't indicate a total lack of suitability.
	ranking := 0
	if !contains(userPrefs.ReadGenders(), candidate.Gender) {
		return -1, nil
	}

	candidateAge := candidate.CalculateAge()
	if candidateAge >= userPrefs.MinAge && candidateAge <= userPrefs.MaxAge {
		// TODO I want to improve this so that I can add the weight of the age gap.
		// so that a smaller gap adds a higher score, and large is a lower score.
		ranking++
	}

	if userPrefs.EnjoysTravel && canPrefs.EnjoysTravel {
		ranking++
	}

	if userPrefs.EducationLevel == canPrefs.EducationLevel {
		ranking++
	}

	if userPrefs.WantsChildren && canPrefs.WantsChildren {
		ranking++
	}

	_, km := haversine.Distance(u.ReadLocation(), candidate.ReadLocation())

	// This ranking based on distance leaves a lot to be desired.. but it gives an idea.
	if km < 1000 {
		ranking += 3
	} else if km < 2000 {
		ranking += 2
	} else if km < 2500 {
		ranking += 1
	}

	return ranking, nil
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
