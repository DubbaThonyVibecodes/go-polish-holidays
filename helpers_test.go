package polishholidays

import (
	"testing"
	"time"
)

func TestGetEasterDate(t *testing.T) {
	for year, want := range easterDates {
		got := getEasterDate(year)
		if got.Format("2006-01-02") != want {
			t.Errorf("getEasterDate(%d) = %s, want %s", year, got.Format("2006-01-02"), want)
		}
		if got.Weekday() != time.Sunday {
			t.Errorf("getEasterDate(%d) = %v, not a Sunday", year, got.Weekday())
		}
	}
}

func TestGetEasterDateRange(t *testing.T) {
	// Easter always falls between 22 March and 25 April inclusive and on a Sunday.
	earliest := false
	latest := false
	for year := 1; year <= 5000; year++ {
		easter := getEasterDate(year)
		if easter.Weekday() != time.Sunday {
			t.Errorf("Easter %d (%v) is not a Sunday", year, easter)
		}
		m, d := easter.Month(), easter.Day()
		if m == time.March && d == 22 {
			earliest = true
		}
		if m == time.April && d == 25 {
			latest = true
		}
		if (m == time.March && d >= 22) || (m == time.April && d <= 25) {
			continue
		}
		t.Errorf("Easter %d = %v, outside 22 March - 25 April", year, easter)
	}
	if !earliest || !latest {
		t.Errorf("expected both extremes to occur in 1..5000 (earliest=%v, latest=%v)", earliest, latest)
	}
}

func TestAddDays(t *testing.T) {
	base := time.Date(2020, time.April, 12, 0, 0, 0, 0, time.UTC)
	cases := []struct {
		days int
		want string
	}{
		{0, "2020-04-12"},
		{1, "2020-04-13"},
		{49, "2020-05-31"},
		{60, "2020-06-11"},
		{-12, "2020-03-31"},
	}
	for _, tc := range cases {
		if got := addDays(base, tc.days).Format("2006-01-02"); got != tc.want {
			t.Errorf("addDays(%v, %d) = %s, want %s", base, tc.days, got, tc.want)
		}
	}
}

func TestIsSameDay(t *testing.T) {
	utc := time.Date(2020, time.April, 12, 0, 0, 0, 0, time.UTC)
	cases := []struct {
		a, b time.Time
		want bool
	}{
		{utc, utc, true},
		{utc, time.Date(2020, time.April, 12, 23, 59, 59, 0, time.UTC), true},
		{utc, time.Date(2020, time.April, 13, 0, 0, 0, 0, time.UTC), false},
		{utc, time.Date(2021, time.April, 12, 0, 0, 0, 0, time.UTC), false},
		{utc, time.Date(2020, time.May, 12, 0, 0, 0, 0, time.UTC), false},
		// Same instant, different calendar day in a different zone: compared by local components.
		{utc, time.Date(2020, time.April, 12, 0, 30, 0, 0, time.FixedZone("UTC+2", 2*3600)), true},
		{utc, time.Date(2020, time.April, 11, 22, 30, 0, 0, time.FixedZone("UTC-2", -2*3600)), false},
	}
	for _, tc := range cases {
		if got := isSameDay(tc.a, tc.b); got != tc.want {
			t.Errorf("isSameDay(%v, %v) = %v, want %v", tc.a, tc.b, got, tc.want)
		}
	}
}
