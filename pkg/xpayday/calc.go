package xpayday

import (
	"time"
)

// SameMonthYear returns true if t1 and t2 hold the same month and year.
func SameMonthYear(t1, t2 time.Time) bool {
	return t1.Year() == t2.Year() && t1.Month() == t2.Month()
}

// FirstPaydayInMonth returns the first pay day for a month, given a reference payday and a pay period.
// panics if period <= 0
func FirstPaydayInMonth(referencePayday time.Time, period time.Duration) time.Time {
	if period <= 0 {
		panic("period <= 0")
	}
	return addUntilMonthChange(referencePayday, -period)
}

// LastPaydayInMonth returns the last pay day for a month, given a reference payday and a pay period.
// panics if period <= 0
func LastPaydayInMonth(referencePayday time.Time, period time.Duration) time.Time {
	if period <= 0 {
		panic("period <= 0")
	}
	return addUntilMonthChange(referencePayday, period)
}

func addUntilMonthChange(referencePayday time.Time, period time.Duration) time.Time {
	if period == 0 {
		panic("period == 0")
	}

	t := referencePayday
	for {
		prev := t.Add(period)
		if !SameMonthYear(prev, referencePayday) {
			return t
		}
		t = prev
	}
}

// PaydaysInMonth returns each payday in the month for given a reference payday and a pay period.
// panics if period <= 0
func PaydaysInMonth(payday time.Time, period time.Duration) (paydays []time.Time) {
	if period <= 0 {
		panic("period <= 0")
	}

	first, last := FirstPaydayInMonth(payday, period), LastPaydayInMonth(payday, period)

	t := first
	for {
		paydays = append(paydays, t)
		t = t.Add(period)
		if t.After(last) {
			break
		}
	}

	return
}
