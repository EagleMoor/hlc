package models

import (
	"time"

	"github.com/bearbin/go-age"
)

type UserGender string

const (
	UserGenderMale   UserGender = "m"
	UserGenderFemale UserGender = "f"
)

//easyjson:json
type User struct {
	ID        int        `json:"id"`
	Email     string     `json:"email"`
	FirstName string     `json:"first_name"`
	LastName  string     `json:"last_name"`
	Gender    UserGender `json:"gender"`
	BirthDate int        `json:"birth_date"`

	visits map[int]*Visit `json:"-"`

	Json []byte `json:"-"`
}

type UserUpdate struct {
	ID        JSONInt    `json:"id"`
	Email     JSONString `json:"email"`
	FirstName JSONString `json:"first_name"`
	LastName  JSONString `json:"last_name"`
	Gender    JSONString `json:"gender"`
	BirthDate JSONInt    `json:"birth_date"`
}

func (uu *UserUpdate) Valid() bool {
	if uu.ID.Set && !uu.ID.Valid {
		return false
	}
	if uu.Email.Set && !uu.Email.Valid {
		return false
	}
	if uu.FirstName.Set && !uu.FirstName.Valid {
		return false
	}
	if uu.LastName.Set && !uu.LastName.Valid {
		return false
	}
	if uu.Gender.Set && !uu.Gender.Valid {
		return false
	}
	if uu.BirthDate.Set && !uu.BirthDate.Valid {
		return false
	}

	return true
}

func (u *User) Valid() bool {
	if u.Gender != UserGenderMale && u.Gender != UserGenderFemale {
		return false
	}

	return true
}

func (u *User) IsMale() bool   { return u.Gender == UserGenderMale }
func (u *User) IsFemale() bool { return u.Gender == UserGenderFemale }

func (u *User) Age() int {
	return age.Age(time.Unix(int64(u.BirthDate), 0))
}
