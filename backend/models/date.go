package models

import (
	"database/sql/driver"
	"encoding/json"
	"strings"
	"time"
)

// Date is a custom type that can handle both date-only and datetime strings
type Date struct {
	time.Time
}

// UnmarshalJSON implements json.Unmarshaler interface
func (d *Date) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	if s == "" || s == "null" {
		d.Time = time.Time{}
		return nil
	}

	// Try parsing as date-only first (YYYY-MM-DD)
	if t, err := time.Parse("2006-01-02", s); err == nil {
		d.Time = t
		return nil
	}

	// Try parsing as datetime (RFC3339)
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		d.Time = t
		return nil
	}

	// Try parsing as datetime without timezone
	if t, err := time.Parse("2006-01-02T15:04:05", s); err == nil {
		d.Time = t
		return nil
	}

	// If all parsing fails, return error
	return json.Unmarshal(b, &d.Time)
}

// MarshalJSON implements json.Marshaler interface
func (d Date) MarshalJSON() ([]byte, error) {
	if d.Time.IsZero() {
		return []byte("null"), nil
	}
	return json.Marshal(d.Time.Format("2006-01-02"))
}

// Value implements driver.Valuer interface for GORM
func (d Date) Value() (driver.Value, error) {
	if d.Time.IsZero() {
		return nil, nil
	}
	return d.Time.Format("2006-01-02"), nil
}

// Scan implements sql.Scanner interface for GORM
func (d *Date) Scan(value interface{}) error {
	if value == nil {
		d.Time = time.Time{}
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		d.Time = v
		return nil
	case string:
		t, err := time.Parse("2006-01-02", v)
		if err != nil {
			return err
		}
		d.Time = t
		return nil
	case []byte:
		t, err := time.Parse("2006-01-02", string(v))
		if err != nil {
			return err
		}
		d.Time = t
		return nil
	default:
		return nil
	}
}
