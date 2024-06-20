package repository

type Swipe struct {
	UserID      int  `json:"user_id,omitempty"`
	CandidateID int  `json:"candidate_id"`
	Likes       bool `json:"likes"`
}
