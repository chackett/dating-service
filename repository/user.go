package repository

type User struct {
	ID          int    `json:"id,omitempty"`
	Email       string `json:"email,omitempty"`
	Password    string `json:"password,omitempty" `
	Name        string `json:"name,omitempty"`
	Gender      string `json:"gender,omitempty"`
	DateOfBirth string `json:"date_of_birth,omitempty"`
}
