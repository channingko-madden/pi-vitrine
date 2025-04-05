package cher

import "time"

// Represent device information
type Device struct {
	Name      string
	Location  string
	CreatedAt time.Time
}
