package models

type EntityType string

const (
	UserEntityType     EntityType = "users"
	LocationEntityType EntityType = "locations"
	VisitsEntityType   EntityType = "visits"
)

func (et EntityType) Invalid() bool {
	return et != UserEntityType && et != LocationEntityType && et != VisitsEntityType
}
