package xpayday

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFirstOfMonth(t *testing.T) {
	t.Run("jan 2023", func(t *testing.T) {
		res := FirstOfMonth(time.Date(2023, time.January, 22, 0, 0, 0, 0, time.UTC))
		assert.Equal(t, time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC), res)
	})
	t.Run("feb 2023", func(t *testing.T) {
		res := FirstOfMonth(time.Date(2023, time.February, 22, 0, 0, 0, 0, time.UTC))
		assert.Equal(t, time.Date(2023, time.February, 1, 0, 0, 0, 0, time.UTC), res)
	})
	t.Run("feb 2024 (leap year)", func(t *testing.T) {
		res := FirstOfMonth(time.Date(2024, time.February, 29, 0, 0, 0, 0, time.UTC))
		assert.Equal(t, time.Date(2024, time.February, 1, 0, 0, 0, 0, time.UTC), res)
	})
	t.Run("wide range", func(t *testing.T) {
		const iterations = 500
		const twentyEightDays = 28 * 24 * time.Hour
		cur := time.Now()
		for i := 0; i < iterations; i++ {
			cur = cur.Add(twentyEightDays)
			res := FirstOfMonth(cur)
			assert.Equal(t, cur.Year(), res.Year())
			assert.Equal(t, cur.Month(), res.Month())
			assert.NotEqual(t, cur.Month(), res.Add(-1).Month())
		}
	})
}

func TestLastOfMonth(t *testing.T) {
	t.Run("jan 2023", func(t *testing.T) {
		res := LastOfMonth(time.Date(2023, time.January, 22, 0, 0, 0, 0, time.UTC))
		assert.Equal(t, time.Date(2023, time.January, 31, 23, 59, 59, 999999999, time.UTC), res)
	})
	t.Run("feb 2023", func(t *testing.T) {
		res := LastOfMonth(time.Date(2023, time.February, 22, 0, 0, 0, 0, time.UTC))
		assert.Equal(t, time.Date(2023, time.February, 28, 23, 59, 59, 999999999, time.UTC), res)
	})
	t.Run("feb 2024 (leap year)", func(t *testing.T) {
		res := LastOfMonth(time.Date(2024, time.February, 22, 0, 0, 0, 0, time.UTC))
		assert.Equal(t, time.Date(2024, time.February, 29, 23, 59, 59, 999999999, time.UTC), res)
	})
	t.Run("wide range", func(t *testing.T) {
		const iterations = 500
		const twentyEightDays = 28 * 24 * time.Hour
		cur := time.Now()
		for i := 0; i < iterations; i++ {
			cur = cur.Add(twentyEightDays)
			res := LastOfMonth(cur)
			assert.Equal(t, cur.Year(), res.Year())
			assert.Equal(t, cur.Month(), res.Month())
			assert.NotEqual(t, cur.Month(), res.Add(1).Month())
		}
	})
}

func TestPaydaysInRange(t *testing.T) {
	const twoWeeks = 2 * 7 * 24 * time.Hour

	t.Run("period must be positive", func(t *testing.T) {
		assert.Panics(t, func() { PaydaysInRange(time.Unix(0, 0), time.Unix(0, 0), time.Unix(1, 0), 0) })
		assert.Panics(t, func() { PaydaysInRange(time.Now(), time.Unix(0, 0), time.Unix(1, 0), -time.Second) })
	})
	t.Run("d1 must less than d2", func(t *testing.T) {
		assert.Panics(t, func() { PaydaysInRange(time.Unix(0, 0), time.Unix(100, 0), time.Unix(100, 0), 1) })
		assert.Panics(t, func() { PaydaysInRange(time.Unix(0, 0), time.Unix(200, 0), time.Unix(100, 0), 1) })
	})

	t.Run("first biweekly payday mar 8th, expect two paydays between mar 5th and mar 23", func(t *testing.T) {
		paydays := []time.Time{
			time.Date(2000, time.March, 8, 0, 0, 0, 0, time.UTC),
			time.Date(2000, time.March, 22, 0, 0, 0, 0, time.UTC),
		}
		res := PaydaysInRange(
			time.Date(2000, time.March, 8, 0, 0, 0, 0, time.UTC),
			time.Date(2000, time.March, 5, 0, 0, 0, 0, time.UTC),
			time.Date(2000, time.March, 23, 0, 0, 0, 0, time.UTC),
			twoWeeks)
		assert.EqualValues(t, paydays, res)
	})

	t.Run("first biweekly payday mar 25th, expect two paydays between apr 5th and apr 23", func(t *testing.T) {
		paydays := []time.Time{
			time.Date(2000, time.April, 8, 0, 0, 0, 0, time.UTC),
			time.Date(2000, time.April, 22, 0, 0, 0, 0, time.UTC),
		}
		res := PaydaysInRange(
			time.Date(2000, time.March, 25, 0, 0, 0, 0, time.UTC),
			time.Date(2000, time.April, 5, 0, 0, 0, 0, time.UTC),
			time.Date(2000, time.April, 23, 0, 0, 0, 0, time.UTC),
			twoWeeks)
		assert.EqualValues(t, paydays, res)
	})
}

func TestPaydaysInMonth(t *testing.T) {
	const twoWeeks = 2 * 7 * 24 * time.Hour

	t.Run("period must be positive", func(t *testing.T) {
		assert.Panics(t, func() { PaydaysInMonth(time.Now(), 0) })
		assert.Panics(t, func() { PaydaysInMonth(time.Now(), -time.Second) })
	})
	t.Run("first biweekly payday july 1st, expect three paydays", func(t *testing.T) {
		paydays := []time.Time{
			time.Date(2000, time.July, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2000, time.July, 15, 0, 0, 0, 0, time.UTC),
			time.Date(2000, time.July, 29, 0, 0, 0, 0, time.UTC),
		}

		for _, refPayday := range paydays {
			v := PaydaysInMonth(refPayday, twoWeeks)
			assert.EqualValues(t, paydays, v)
		}
	})
	t.Run("first biweekly payday july 4th, expect two paydays", func(t *testing.T) {
		paydays := []time.Time{
			time.Date(2000, time.July, 4, 0, 0, 0, 0, time.UTC),
			time.Date(2000, time.July, 18, 0, 0, 0, 0, time.UTC),
		}

		for _, refPayday := range paydays {
			v := PaydaysInMonth(refPayday, twoWeeks)
			assert.EqualValues(t, paydays, v)
		}
	})

	t.Run("first and only payday july 4th", func(t *testing.T) {
		paydays := []time.Time{
			time.Date(2000, time.July, 4, 0, 0, 0, 0, time.UTC),
		}

		for _, refPayday := range paydays {
			v := PaydaysInMonth(refPayday, 24*100*time.Hour)
			assert.EqualValues(t, paydays, v)
		}
	})
}
