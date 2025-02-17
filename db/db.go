package db

import "time"

func parseCreatedTime(s string) time.Time {
	time, err := time.Parse(time.RFC3339Nano, s)
	if err != nil {
		panic(err)
	}

	return time
}
