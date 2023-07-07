package xpayday

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSameMonthYear(t *testing.T) {
	type args struct {
		t1 time.Time
		t2 time.Time
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "same month",
			args: args{
				t1: time.Date(2000 /*year*/, 2 /*mon*/, 2 /*day*/, 0 /*hour*/, 0 /*min*/, 0 /*sec*/, 0 /*nsec*/, time.UTC),
				t2: time.Date(2000 /*year*/, 2 /*mon*/, 24 /*day*/, 0 /*hour*/, 0 /*min*/, 0 /*sec*/, 0 /*nsec*/, time.UTC),
			},
			want: true,
		},
		{
			name: "same month, different year",
			args: args{
				t1: time.Date(2000 /*year*/, time.February /*mon*/, 2 /*day*/, 0 /*hour*/, 0 /*min*/, 0 /*sec*/, 0 /*nsec*/, time.UTC),
				t2: time.Date(2001 /*year*/, time.February /*mon*/, 24 /*day*/, 0 /*hour*/, 0 /*min*/, 0 /*sec*/, 0 /*nsec*/, time.UTC),
			},
			want: false,
		},
		{
			name: "same month via normalization (days: 40th February == March)",
			args: args{
				t1: time.Date(2000 /*year*/, time.March /*mon*/, 2 /*day*/, 0 /*hour*/, 0 /*min*/, 0 /*sec*/, 0 /*nsec*/, time.UTC),
				t2: time.Date(2000 /*year*/, time.February /*mon*/, 40 /*day*/, 0 /*hour*/, 0 /*min*/, 0 /*sec*/, 0 /*nsec*/, time.UTC),
			},
			want: true,
		},
		{
			name: "different month via normalization (days: 31st February != February)",
			args: args{
				t1: time.Date(2000 /*year*/, time.February /*mon*/, 2 /*day*/, 0 /*hour*/, 0 /*min*/, 0 /*sec*/, 0 /*nsec*/, time.UTC),
				t2: time.Date(2000 /*year*/, time.February /*mon*/, 31 /*day*/, 0 /*hour*/, 0 /*min*/, 0 /*sec*/, 0 /*nsec*/, time.UTC),
			},
			want: false,
		},
		{
			name: "different month via normalization (hours: February 2 744:15:15 != February)",
			args: args{
				t1: time.Date(2000 /*year*/, time.February /*mon*/, 2 /*day*/, 0 /*hour*/, 0 /*min*/, 0 /*sec*/, 0 /*nsec*/, time.UTC),
				t2: time.Date(2000 /*year*/, time.February /*mon*/, 2 /*day*/, 24*31 /*hour*/, 15 /*min*/, 15 /*sec*/, 0 /*nsec*/, time.UTC),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SameMonthYear(tt.args.t1, tt.args.t2); got != tt.want {
				t.Errorf("SameMonth() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFirstPaydayInMonth(t *testing.T) {
	const oneweek = 7 * 24 * time.Hour

	t.Run("period must be positive", func(t *testing.T) {
		assert.Panics(t, func() { FirstPaydayInMonth(time.Now(), 0) })
		assert.Panics(t, func() { FirstPaydayInMonth(time.Now(), -time.Second) })
	})
	t.Run("reference date is first, result is reference date", func(t *testing.T) {
		april3rd := time.Date(2000, time.April, 3, 0, 0, 0, 0, time.UTC)
		v := FirstPaydayInMonth(april3rd, oneweek)
		assert.Equal(t, april3rd, v)
	})
	t.Run("reference date is last, result is earlier", func(t *testing.T) {
		april3rd := time.Date(2000, time.April, 3, 0, 0, 0, 0, time.UTC)
		april24th := time.Date(2000, time.April, 24, 0, 0, 0, 0, time.UTC)
		v := FirstPaydayInMonth(april24th, oneweek)
		assert.Equal(t, april3rd, v)
	})
	t.Run("interval > 1 month, result is reference date", func(t *testing.T) {
		april24th := time.Date(2000, time.April, 24, 0, 0, 0, 0, time.UTC)
		v := FirstPaydayInMonth(april24th, 5*oneweek)
		assert.Equal(t, april24th, v)
	})
}

func TestLastPaydayInMonth(t *testing.T) {
	const oneweek = 7 * 24 * time.Hour

	t.Run("period must be positive", func(t *testing.T) {
		assert.Panics(t, func() { LastPaydayInMonth(time.Now(), 0) })
		assert.Panics(t, func() { LastPaydayInMonth(time.Now(), -time.Second) })
	})
	t.Run("reference date is last, result is reference date", func(t *testing.T) {
		april24th := time.Date(2000, time.April, 24, 0, 0, 0, 0, time.UTC)
		v := LastPaydayInMonth(april24th, oneweek)
		assert.Equal(t, april24th, v)
	})
	t.Run("reference date is first, result is later", func(t *testing.T) {
		april3rd := time.Date(2000, time.April, 3, 0, 0, 0, 0, time.UTC)
		april24th := time.Date(2000, time.April, 24, 0, 0, 0, 0, time.UTC)
		v := LastPaydayInMonth(april3rd, oneweek)
		assert.Equal(t, april24th, v)
	})
	t.Run("interval > 1 month, result is reference date", func(t *testing.T) {
		april24th := time.Date(2000, time.April, 24, 0, 0, 0, 0, time.UTC)
		v := LastPaydayInMonth(april24th, 5*oneweek)
		assert.Equal(t, april24th, v)
	})
}

func TestPaydaysInMonth(t *testing.T) {
	const biweekly = 2 * 7 * 24 * time.Hour
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
			v := PaydaysInMonth(refPayday, biweekly)
			assert.EqualValues(t, paydays, v)
		}
	})
	t.Run("first biweekly payday july 4th, expect two paydays", func(t *testing.T) {
		paydays := []time.Time{
			time.Date(2000, time.July, 4, 0, 0, 0, 0, time.UTC),
			time.Date(2000, time.July, 18, 0, 0, 0, 0, time.UTC),
		}

		for _, refPayday := range paydays {
			v := PaydaysInMonth(refPayday, biweekly)
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
