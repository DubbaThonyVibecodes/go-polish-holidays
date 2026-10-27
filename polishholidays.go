// Package polishholidays computes the dates of public holidays in Poland.
//
// It is a Go port of the poland-public-holidays JavaScript library by
// Kamil Szydłowski (https://github.com/szydlovski/poland-public-holidays).
//
// All returned dates are at midnight UTC. Input dates are compared by their
// calendar day in their own location, so a time.Time in any zone works as
// expected.
package polishholidays

import (
	"slices"
	"time"
)

// GetHolidaysInYear returns every public holiday in the given year, sorted by
// date. It returns an error wrapping ErrInvalidYear if year is below 1.
func GetHolidaysInYear(year int) ([]Holiday, error) {
	if err := validateYear(year); err != nil {
		return nil, err
	}

	easterDate := getEasterDate(year)
	holidays := make([]Holiday, 0, len(holidayRules))
	for _, rule := range holidayRules {
		var date time.Time
		switch rule.typ {
		case fixed:
			date = time.Date(year, rule.month, rule.day, 0, 0, 0, 0, time.UTC)
		case movable:
			date = addDays(easterDate, rule.afterEaster)
		}
		holidays = append(holidays, Holiday{
			Name:   rule.name,
			NamePL: rule.namePL,
			Date:   date,
		})
	}

	slices.SortFunc(holidays, func(a, b Holiday) int {
		return a.Date.Compare(b.Date)
	})
	return holidays, nil
}

// GetHolidayOnDate returns the public holiday observed on the given date and
// true, or a zero Holiday and false if the date is not a public holiday.
// Dates before year 1 are never holidays.
func GetHolidayOnDate(date time.Time) (Holiday, bool) {
	holidays, err := GetHolidaysInYear(date.Year())
	if err != nil {
		return Holiday{}, false
	}
	for _, holiday := range holidays {
		if isSameDay(holiday.Date, date) {
			return holiday, true
		}
	}
	return Holiday{}, false
}

// IsHoliday reports whether there is a public holiday on the given date.
func IsHoliday(date time.Time) bool {
	_, ok := GetHolidayOnDate(date)
	return ok
}

// IsWeekend reports whether the given date falls on a Saturday or Sunday.
func IsWeekend(date time.Time) bool {
	weekday := date.Weekday()
	return weekday == time.Saturday || weekday == time.Sunday
}

// IsWeekendOrHoliday reports whether the given date is a weekend day or a
// public holiday.
func IsWeekendOrHoliday(date time.Time) bool {
	return IsWeekend(date) || IsHoliday(date)
}

// IsWorkingDay reports whether the given date is neither a weekend day nor a
// public holiday.
func IsWorkingDay(date time.Time) bool {
	return !IsWeekendOrHoliday(date)
}
