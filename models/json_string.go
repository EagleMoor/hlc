package models

import "encoding/json"

// JSONString for string with "null" value
type JSONString struct {
	Value string
	Valid bool
	Set   bool
}

// UnmarshalJSON string with "null" value
func (i *JSONString) UnmarshalJSON(data []byte) error {
	// If this method was called, the value was set.
	i.Set = true

	if string(data) == "null" {
		// The key was set to null
		i.Valid = false
		return nil
	}

	// The key isn't set to null
	var temp string
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}
	i.Value = temp
	i.Valid = true
	return nil
}
