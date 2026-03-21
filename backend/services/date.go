package services

import (
	"encoding/json"
	"strings"
	"time"
)

// FlexDate parses JSON date strings in "2006-01-02" or RFC3339 format (e.g. from HTML date input).
type FlexDate struct {
	time.Time
}

func (d *FlexDate) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	if s == "" || s == "null" {
		d.Time = time.Time{}
		return nil
	}
	// Try date-only first (what frontend sends from <input type="date">)
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		t, err = time.Parse(time.RFC3339, s)
		if err != nil {
			return err
		}
	}
	d.Time = t
	return nil
}

func (d FlexDate) MarshalJSON() ([]byte, error) {
	if d.Time.IsZero() {
		return []byte(`null`), nil
	}
	return json.Marshal(d.Time.Format("2006-01-02"))
}
