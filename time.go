package notion

import (
	"encoding/json"
	"errors"
	"time"
)

// Length of a date string, e.g. `2006-01-02`.
const dateLength = 10

// DateTimeFormat is used when calling time.Parse, using RFC3339 with microsecond
// precision, which is what the Notion API returns in JSON response data.
const DateTimeFormat = "2006-01-02T15:04:05.999Z07:00"

// DateTime represents a Notion date property with optional time.
type DateTime struct {
	time.Time
	hasTime bool
}

// ParseDateTime parses an RFC3339 formatted string with optional time.
func ParseDateTime(value string) (DateTime, error) {
	if len(value) > len(DateTimeFormat) {
		return DateTime{}, errors.New("invalid datetime string")
	}

	t, err := time.Parse(DateTimeFormat[:len(value)], value)
	if err != nil {
		return DateTime{}, err
	}

	dt := DateTime{
		Time:    t,
		hasTime: len(value) > dateLength,
	}

	return dt, nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (dt *DateTime) UnmarshalJSON(b []byte) error {
	if len(b) < 2 {
		return errors.New("invalid datetime string")
	}

	parsed, err := ParseDateTime(string(b[1 : len(b)-1]))
	if err != nil {
		return err
	}

	*dt = parsed

	return nil
}

// MarshalJSON implements json.Marshaler. It returns an RFC399 formatted string,
// using microsecond precision ()
func (dt DateTime) MarshalJSON() ([]byte, error) {
	if dt.hasTime {
		return json.Marshal(dt.Time)
	}
	return []byte(`"` + dt.Time.Format(DateTimeFormat[:dateLength]) + `"`), nil
}

// NewDateTime returns a new DateTime. If `haseTime` is true, time is included
// when encoding to JSON.
func NewDateTime(t time.Time, hasTime bool) DateTime {
	var tt time.Time

	if hasTime {
		tt = t
	} else {
		tt = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
	}

	return DateTime{
		Time:    tt,
		hasTime: hasTime,
	}
}

// HasTime returns true if the datetime was parsed from a string that included time.
func (dt *DateTime) HasTime() bool {
	return dt.hasTime
}

// Equal returns true if both DateTime values have equal underlying time.Time and
// hasTime fields.
func (dt DateTime) Equal(value DateTime) bool {
	if !dt.Time.Equal(value.Time) {
		return false
	}
	return dt.hasTime == value.hasTime
}
