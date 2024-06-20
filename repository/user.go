package repository

import (
	"time"
)

type User struct {
	ID          int        `json:"id,omitempty"`
	Email       string     `json:"email,omitempty"`
	Password    string     `json:"password,omitempty" `
	Name        string     `json:"name,omitempty"`
	Gender      string     `json:"gender,omitempty"`
	DateOfBirth *time.Time `json:"date_of_birth,omitempty"`
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

func (u *User) RankCandidate(candidate User) (int, error) {
	return 0, nil
}
