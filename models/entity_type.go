package models

// EntityType type of entity
type EntityType string

const (
	// UserEntityType user
	UserEntityType EntityType = "users"

	// LocationEntityType locations
	LocationEntityType EntityType = "locations"

	// VisitsEntityType visits
	VisitsEntityType EntityType = "visits"
)

// Invalid Entity Type
func (et EntityType) Invalid() bool {
	return et != UserEntityType && et != LocationEntityType && et != VisitsEntityType
}
