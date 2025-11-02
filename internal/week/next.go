package week

import "time"

func NextMonday(from time.Time) time.Time {
	loc := from.Location()
	daysToMonday := (int(time.Monday) - int(from.Weekday()) + 7) % 7
	if from.Weekday() == time.Sunday {
		daysToMonday = 1
	}
	return time.Date(from.Year(), from.Month(), from.Day(), 0, 0, 0, 0, loc).AddDate(0, 0, daysToMonday)
}
