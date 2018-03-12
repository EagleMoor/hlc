package models

// Visit user
//easyjson:json
type Visit struct {
	ID         int `json:"id"`
	LocationID int `json:"location"`
	UserID     int `json:"user"`
	VisitedAt  int `json:"visited_at"`
	Mark       int `json:"mark"`

	User *User `json:"-"`

	JSON []byte `json:"-"`
}

// VisitUpdate JSON for update Visit
type VisitUpdate struct {
	ID         JSONInt `json:"id"`
	LocationID JSONInt `json:"location"`
	UserID     JSONInt `json:"user"`
	VisitedAt  JSONInt `json:"visited_at"`
	Mark       JSONInt `json:"mark"`
}

// Valid JSON for update
func (vu *VisitUpdate) Valid() bool {
	if vu.ID.Set && !vu.ID.Valid {
		return false
	}
	if vu.LocationID.Set && !vu.LocationID.Valid {
		return false
	}
	if vu.UserID.Set && !vu.UserID.Valid {
		return false
	}
	if vu.VisitedAt.Set && !vu.VisitedAt.Valid {
		return false
	}
	if vu.Mark.Set && !vu.Mark.Valid {
		return false
	}

	return true
}

// Valid Visit
func (v *Visit) Valid() bool {
	return true
}
