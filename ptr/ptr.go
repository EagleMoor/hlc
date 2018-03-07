package ptr

// String returns pointer to a new variable for specified string value
func String(s string) *string {
	return &s
}

// Int returns pointer to a new variable for specified int value
func Int(n int) *int {
	return &n
}

// Int64 returns pointer to a new variable for specified int64 value
func Int64(n int64) *int64 {
	return &n
}

// Float64 returns pointer to a new variable for specified float64 value
func Float64(n float64) *float64 {
	return &n
}

// Bool returns pointer to a new variable for specified bool value
func Bool(b bool) *bool {
	return &b
}
