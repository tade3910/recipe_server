package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type StringList []string

// Scan implements the sql.Scanner interface for StringList
func (s *StringList) Scan(src any) error {
	bytes, ok := src.([]byte)
	if !ok {
		return errors.New("src value cannot cast to []byte")
	}
	err := json.Unmarshal(bytes, &s)
	if err != nil {
		return err
	}
	return nil
}

// Value implements the driver.Valuer interface for StringList
func (s StringList) Value() (driver.Value, error) {
	if len(s) == 0 {
		return nil, nil
	}
	// Convert to JSON before storing
	return json.Marshal(s)
}
