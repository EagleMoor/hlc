package models

//easyjson:json

// Location of user
type Location struct {
	ID       int    `json:"id"`
	Place    string `json:"place"`
	Country  string `json:"country"`
	City     string `json:"city"`
	Distance int    `json:"distance"`

	visits map[int]*Visit

	JSON []byte `json:"-"`
}

// LocationUpdate for update Location
type LocationUpdate struct {
	ID       JSONInt    `json:"id"`
	Place    JSONString `json:"place"`
	Country  JSONString `json:"country"`
	City     JSONString `json:"city"`
	Distance JSONInt    `json:"distance"`
}

// Valid check
func (lu *LocationUpdate) Valid() bool {
	if lu.ID.Set && !lu.ID.Valid {
		return false
	}

	if lu.Place.Set && !lu.Place.Valid {
		return false
	}

	if lu.Country.Set && !lu.Country.Valid {
		return false
	}

	if lu.City.Set && !lu.City.Valid {
		return false
	}

	if lu.Distance.Set && !lu.Distance.Valid {
		return false
	}

	return true
}

// Valid check
func (l *Location) Valid() bool {
	return true
}

// GetAVG for user
func (l *Location) GetAVG(fromDate, toDate, fromAge, toAge *int, gender *UserGender) float32 {
	var count, sum float32

	if l.visits == nil {
		return 0.0
	}

	for _, visit := range l.visits {
		if fromDate != nil && visit.VisitedAt <= *fromDate {
			continue
		}

		if toDate != nil && visit.VisitedAt >= *toDate {
			continue
		}

		if fromAge != nil {
			if visit.User.Age() < *fromAge {
				continue
			}
		}

		if toAge != nil {
			if visit.User.Age() >= *toAge {
				continue
			}
		}

		if gender != nil && visit.User.Gender != *gender {
			continue
		}

		count++
		sum += float32(visit.Mark)
	}

	if count == 0 {
		return 0.0
	}

	return sum / count
}
