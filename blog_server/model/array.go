package model

import (
	"database/sql/driver"
	"strings"
)

type Array []string

func (m *Array) Scan(val interface{}) error {
	s := val.([]uint8)
	ss := strings.Split(string(s), "|")
	*m = ss
	return nil
}

func (m Array) Value() (driver.Value, error) {
	str := strings.Join(m, "|")
	return str, nil
}
