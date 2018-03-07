package models

//easyjson:json
type Visit struct {
	ID         int `json:"id"`
	LocationID int `json:"location"`
	UserID     int `json:"user"`
	VisitedAt  int `json:"visited_at"`
	Mark       int `json:"mark"`

	User *User `json:"-"`

	Json []byte `json:"-"`
}

type VisitUpdate struct {
	ID         JSONInt `json:"id"`
	LocationID JSONInt `json:"location"`
	UserID     JSONInt `json:"user"`
	VisitedAt  JSONInt `json:"visited_at"`
	Mark       JSONInt `json:"mark"`
}

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

func (v *Visit) Valid() bool {
	return true
}
