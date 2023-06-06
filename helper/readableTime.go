package helper

import "time"

func ReadableTime(i string) time.Time {
	FORMAT := "02 Jan 2006"
	t, _ := time.Parse(FORMAT, i)

	return t
}
