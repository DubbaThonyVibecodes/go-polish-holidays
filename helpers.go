package polishholidays

import (
	"errors"
	"fmt"
	"time"
)

var (
	// ErrInvalidDate is returned by ParseDate when the input cannot be
	// interpreted as a date.
	ErrInvalidDate = errors.New("invalid date")
	// ErrInvalidYear is returned by GetHolidaysInYear for years below 1.
	ErrInvalidYear = errors.New("invalid year")
)

// dateLayouts lists the string formats accepted by ParseDate, in the order
// they are tried.
var dateLayouts = []string{
	"2006-01-02",
	time.RFC3339,
	"2006-01-02T15:04:05",
	"2006-01-02T15:04",
	"2006-01-02 15:04:05",
}

// addDays returns date shifted by the given number of calendar days.
func addDays(date time.Time, days int) time.Time {
	return date.AddDate(0, 0, days)
}

// isSameDay reports whether both times fall on the same calendar day, each
// evaluated in its own location.
func isSameDay(date1, date2 time.Time) bool {
	y1, m1, d1 := date1.Date()
	y2, m2, d2 := date2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

// getEasterDate computes the date of Easter Sunday in the Gregorian calendar
// for the given year, at midnight UTC.
//
// https://en.wikipedia.org/wiki/Date_of_Easter#Anonymous_Gregorian_algorithm
func getEasterDate(year int) time.Time {
	a := year % 19
	b := year / 100
	c := year % 100
	d := b / 4
	e := b % 4
	f := (b + 8) / 25
	g := (b - f + 1) / 3
	h := (19*a + b - d - g + 15) % 30
	i := c / 4
	k := c % 4
	l := (32 + 2*e + 2*i - h - k) % 7
	m := (a + 11*h + 22*l) / 451
	p := h + l - 7*m + 114

	day := p%31 + 1
	month := time.Month(p / 31)

	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}

// validateYear returns ErrInvalidYear for years below 1.
func validateYear(year int) error {
	if year < 1 {
		return fmt.Errorf("%w: %d", ErrInvalidYear, year)
	}
	return nil
}

// ParseDate parses an ISO 8601 date ("2006-01-02") or date-time (RFC 3339,
// or a date-time without a zone offset) string. Values without an explicit
// zone are interpreted as UTC. It returns an error wrapping ErrInvalidDate if
// the string does not match any supported layout.
//
// It is the Go counterpart of the JavaScript library's acceptance of date
// strings; combine it with IsHoliday or GetHolidayOnDate as needed.
func ParseDate(s string) (time.Time, error) {
	for _, layout := range dateLayouts {
		if t, err := time.Parse(layout, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("%w: %q", ErrInvalidDate, s)
}
