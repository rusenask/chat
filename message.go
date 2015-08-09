package main

import (
	"time"
)

// message represents a single msg
// now attaching sender's name, message and time
type message struct {
	Name    string
	Message string
	When    time.Time
}
