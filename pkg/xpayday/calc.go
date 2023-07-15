package xpayday

import (
	"time"
)

// FirstOfMonth returns the first time.Time value for the month in the given ref month.
func FirstOfMonth(ref time.Time) time.Time {
	return time.Date(ref.Year(), ref.Month(), 1, 0, 0, 0, 0, ref.Location())
}

// FirstOfMonth returns the last time.Time value for the month in the given ref month.
func LastOfMonth(ref time.Time) time.Time {
	return time.Date(ref.Year(), ref.Month()+1, 1, 0, 0, 0, 0, ref.Location()).Add(-1)
}

// PaydaysInRange returns each payday in the date rage [d1, d2) for given a reference payday and a pay period.
// panics if period <= 0 or d1 >= d2
func PaydaysInRange(payday, d1, d2 time.Time, period time.Duration) (paydays []time.Time) {
	if period <= 0 {
		panic("period <= 0")
	}
	if d1.Compare(d2) != -1 {
		panic("d1 >= d2")
	}

	var first time.Time
	switch payday.Compare(d1) {
	case -1: // payday < d1
		diff := d1.Sub(payday)
		truncDiff := diff.Abs().Truncate(period)
		first = payday.Add(truncDiff)
		if first.Before(d1) {
			first = first.Add(period)
		}
	case 1: //payday > d1
		diff := d1.Sub(payday)
		truncDiff := diff.Abs().Truncate(period)
		first = payday.Add(-truncDiff)
		if first.Before(d1) {
			first = first.Add(period)
		}
	default:
		first = d1
	}

	last := d2
	t := first
	for {
		paydays = append(paydays, t)
		t = t.Add(period)
		if t.Compare(last) != -1 {
			break
		}
	}

	return
}

// PaydaysInMonth returns each payday in the month for given a reference payday and a pay period.
// panics if period <= 0
func PaydaysInMonth(payday time.Time, period time.Duration) (paydays []time.Time) {
	if period <= 0 {
		panic("period <= 0")
	}

	return PaydaysInRange(payday, FirstOfMonth(payday), LastOfMonth(payday), period)
}
