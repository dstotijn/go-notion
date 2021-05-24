package notion

import "time"

// StringPtr returns the pointer of a string value.
func StringPtr(s string) *string {
	return &s
}

// IntPtr returns the pointer of an int value.
func IntPtr(i int) *int {
	return &i
}

// BoolPtr returns the pointer of a bool value.
func BoolPtr(b bool) *bool {
	return &b
}

// TimePtr returns the pointer of a time.Time value.
func TimePtr(t time.Time) *time.Time {
	return &t
}

// Float64Ptr returns the pointer of a float64 value.
func Float64Ptr(f float64) *float64 {
	return &f
}
