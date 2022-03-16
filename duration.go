package config

import (
	"errors"
	"fmt"
	"time"
)

// Duration wraps a time.Duration with the ability to marshal and unmarshal to JSON the same as it
// marshals to text. For some reason time.Duration text marshaller uses "5s", but JSON uses
// nanoseconds as an integer. For configs we want environment and JSON configs to use the same
// values.
//
// The main disadvantages are that you must use NewDuration to create them from a time.Duration and
// you must add .Duration to use it as a time.Duration.
type Duration struct {
	time.Duration
}

func NewDuration(d time.Duration) Duration {
	return Duration{d}
}

func (v Duration) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", v.String())), nil
}

func (v *Duration) UnmarshalJSON(js []byte) error {
	l := len(js)
	if l < 2 {
		return errors.New("Too short")
	}
	if js[0] != '"' || js[l-1] != '"' {
		return errors.New("Too short")
	}

	duration, err := time.ParseDuration(string(js[1 : l-1]))
	if err != nil {
		return err
	}

	*v = NewDuration(duration)
	return nil
}

func (v Duration) MarshalText() ([]byte, error) {
	return []byte(v.String()), nil
}

func (v *Duration) UnmarshalText(text []byte) error {
	duration, err := time.ParseDuration(string(text))
	if err != nil {
		return err
	}

	*v = NewDuration(duration)
	return nil
}
