# go-polish-holidays

[![Go Reference](https://pkg.go.dev/badge/github.com/DubbaThonyVibecodes/go-polish-holidays.svg)](https://pkg.go.dev/github.com/DubbaThonyVibecodes/go-polish-holidays)

A small but comprehensive, fully tested Go library that computes the dates of public holidays in Poland.

This is a 1:1 port of [poland-public-holidays](https://github.com/szydlovski/poland-public-holidays) by Kamil Szydłowski. It has no dependencies outside the standard library.

Port performed by Claude Fable 5.1

## Usage

```
go get github.com/DubbaThonyVibecodes/go-polish-holidays
```

```go
import polishholidays "github.com/DubbaThonyVibecodes/go-polish-holidays"

// accepts time.Time values
polishholidays.IsHoliday(time.Date(2018, time.May, 1, 0, 0, 0, 0, time.UTC)) // true - Labour Day, a fixed holiday

// or parse ISO 8601 / RFC 3339 strings with ParseDate
easter, _ := polishholidays.ParseDate("2020-04-12")
polishholidays.IsHoliday(easter) // true - Easter, a movable holiday

date, _ := polishholidays.ParseDate("2016-05-26")
holiday, ok := polishholidays.GetHolidayOnDate(date)
holiday.Name // "Feast of Corpus Christi"
ok           // true

date, _ = polishholidays.ParseDate("1978-10-02")
_, ok = polishholidays.GetHolidayOnDate(date)
ok // false - not a public holiday

holidaysIn2020, _ := polishholidays.GetHolidaysInYear(2020)
len(holidaysIn2020)        // 13
holidaysIn2020[0].Name     // "New Year"
holidaysIn2020[0].NamePL   // "Nowy Rok"
holidaysIn2020[0].Date     // 2020-01-01 00:00:00 +0000 UTC
```

## API

### `GetHolidaysInYear(year int) ([]Holiday, error)`

Returns a slice of all public holidays in the given year, sorted by date. Returns an error wrapping `ErrInvalidYear` if `year` is below 1. Each `Holiday` has the following fields:

- `Date` - a `time.Time` set to midnight UTC on the day the holiday is observed
- `Name` - the English name of the holiday
- `NamePL` - the Polish name of the holiday

### `GetHolidayOnDate(date time.Time) (Holiday, bool)`

Returns the holiday observed on the given date and `true`, or a zero `Holiday` and `false` if there is no public holiday on that date.

### `IsHoliday(date time.Time) bool`

Returns `true` if there is a public holiday on the given date, and `false` if there is not.

### `IsWeekend(date time.Time) bool`, `IsWeekendOrHoliday(date time.Time) bool`, `IsWorkingDay(date time.Time) bool`

Convenience helpers. A working day is any day that is neither a Saturday, a Sunday, nor a public holiday.

### `ParseDate(s string) (time.Time, error)`

Parses an ISO 8601 date (`2006-01-02`) or an RFC 3339 date-time. Date-times without a zone offset are interpreted as UTC. Returns an error wrapping `ErrInvalidDate` if the string is not a recognised date.

## Notes on the port

- The JavaScript version accepts `Date` objects, date strings or numbers in every function. Go is statically typed, so dates are `time.Time`, years are `int`, and string input goes through `ParseDate`.
- Dates are compared by calendar day in the `time.Time`'s own location, so a value in any zone is treated as the day it reads as in that zone.
- Where the JavaScript version throws, the Go version returns an error (`GetHolidaysInYear`, `ParseDate`) or `false` (`GetHolidayOnDate` for dates before year 1).
- `testdata/holidays_1900_2100.json` was generated with the original library and the test suite verifies the port against it for every year from 1900 to 2100.

## License

MIT, inherited unchanged from the original library. See [LICENSE](LICENSE).
