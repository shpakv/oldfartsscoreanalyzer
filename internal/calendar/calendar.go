package calendar

import (
	"context"
	"oldfartscounter/internal/week"
	"time"
)

type (
	Calendar interface {
		NextWeekGameDays(ctx context.Context, from time.Time) ([]Day, error)
	}

	Day struct {
		Type DayType
		Date time.Time
	}
	DayType int
)

const (
	PlayDay DayType = iota
	ThinkingDay
)

type defaultCalendar struct{}

func NewDefaultCalendar() Calendar {
	return &defaultCalendar{}
}

func (d *defaultCalendar) NextWeekGameDays(_ context.Context, from time.Time) ([]Day, error) {
	nextMonday := week.NextMonday(from)
	fridayDate := nextMonday.AddDate(0, 0, 4)
	saturdayDate := nextMonday.AddDate(0, 0, 5)
	thursdayDate := nextMonday.AddDate(0, 0, 3)

	return []Day{
		{
			Type: ThinkingDay,
			Date: thursdayDate,
		},
		{
			Type: PlayDay,
			Date: fridayDate,
		},
		{
			Type: PlayDay,
			Date: saturdayDate,
		},
	}, nil
}
