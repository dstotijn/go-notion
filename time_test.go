package notion_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/dstotijn/go-notion"
	"github.com/google/go-cmp/cmp"
)

func mustParseDateTime(value string) notion.DateTime {
	dt, err := notion.ParseDateTime(value)
	if err != nil {
		panic(err)
	}
	return dt
}

func TestTimeMarshalJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		dateTime notion.DateTime
		expJSON  []byte
	}{
		{
			name:     "date and time",
			dateTime: mustParseDateTime("2021-05-23T09:11:50.123Z"),
			expJSON:  []byte(`"2021-05-23T09:11:50.123Z"`),
		},
		{
			name:     "date without time",
			dateTime: mustParseDateTime("2021-05-23"),
			expJSON:  []byte(`"2021-05-23"`),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			dtJSON, err := json.Marshal(tt.dateTime)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if diff := cmp.Diff(string(tt.expJSON), string(dtJSON)); diff != "" {
				t.Fatalf("encoded JSON not equal (-exp, +got):\n%v", diff)
			}

		})
	}
}

func TestTimeUnmarshalJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		timeString  string
		expDateTime notion.DateTime
		expHasTime  bool
		expError    error
	}{
		{
			name:        "date and time",
			timeString:  "2021-05-23T09:11:50.123+00:00",
			expDateTime: notion.NewDateTime(mustParseTime(time.RFC3339Nano, "2021-05-23T09:11:50.123Z"), true),
			expHasTime:  true,
			expError:    nil,
		},
		{
			name:        "date without time",
			timeString:  "2021-05-23",
			expDateTime: notion.NewDateTime(mustParseTime(time.RFC3339Nano, "2021-05-23T09:11:50.123Z"), false),
			expHasTime:  false,
			expError:    nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			type testDateTime struct {
				DateTime notion.DateTime `json:"time"`
			}

			var dt testDateTime
			err := json.Unmarshal([]byte(`{"time":"`+tt.timeString+`"}`), &dt)

			if tt.expError == nil && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.expError != nil && err == nil {
				t.Fatalf("error not equal (expected: %v, got: nil)", tt.expError)
			}
			if tt.expError != nil && err != nil && tt.expError.Error() != err.Error() {
				t.Fatalf("error not equal (expected: %v, got: %v)", tt.expError, err)
			}

			if diff := cmp.Diff(tt.expDateTime.Time, dt.DateTime.Time); diff != "" {
				t.Fatalf("time not equal (-exp, +got):\n%v", diff)
			}

			if tt.expHasTime != dt.DateTime.HasTime() {
				t.Fatalf("has time not equal (expected: %v, got: %v)", tt.expHasTime, dt.DateTime.HasTime())
			}
		})
	}
}
