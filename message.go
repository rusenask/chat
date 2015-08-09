package main

import (
	"fmt"
	"time"
)

// Marshaler implements Unmarshaler interface
type Marshaler interface {
	MarshalJSON() ([]byte, error)
}

// JSONTime formatting time
type JSONTime time.Time

// MarshalJSON wraps time.Time to return formated version
func (t JSONTime) MarshalJSON() ([]byte, error) {
	//do your serializing here
	stamp := fmt.Sprintf("\"%s\"", time.Time(t).Format("Mon Jan _2 15:04:05"))
	return []byte(stamp), nil
}

// message represents a single msg
// now attaching sender's name, message and time
type message struct {
	Name    string
	Message string
	When    JSONTime
}
