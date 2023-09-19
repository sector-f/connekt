package main

import "time"

type eventMessage struct {
	Station string
	Title   string
	Artist  string
	Album   string
	Comment string

	timestamp time.Time
}
