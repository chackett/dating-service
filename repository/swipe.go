package repository

type Swipe struct {
	UserID      int  `json:"userId,omitempty"`
	CandidateID int  `json:"candidateId"`
	Likes       bool `json:"likes"`
}
