package repository

import "strings"

type UserPreferences struct {
	UserID         int    `json:"user_id"`
	WantsChildren  bool   `json:"wants_children"`
	EnjoysTravel   bool   `json:"enjoys_travel"`
	EducationLevel string `json:"education_level"`
	MinAge         int    `json:"min_age"`
	MaxAge         int    `json:"max_age"`
	Genders        string `json:"genders"`
}

func (u *UserPreferences) ReadGenders() []string {
	return strings.Split(u.Genders, ",")
}
