package main

import (
	"time"
)

// message represents a single msg
type message struct {
	Name    string
	Message string
	When    time.Time
}
