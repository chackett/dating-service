package repository

import "strings"

type UserPreferences struct {
	UserID         int    `json:"userId"`
	WantsChildren  bool   `json:"wantsChildren"`
	EnjoysTravel   bool   `json:"enjoysTravel"`
	EducationLevel string `json:"educationLevel"`
	MinAge         int    `json:"minAge"`
	MaxAge         int    `json:"maxAge"`
	Genders        string `json:"genders"`
}

func (u *UserPreferences) ReadGenders() []string {
	return strings.Split(u.Genders, ",")
}
