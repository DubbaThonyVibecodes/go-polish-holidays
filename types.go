package polishholidays

import "time"

// Holiday describes a single public holiday in a given year.
type Holiday struct {
	// Name is the English name of the holiday.
	Name string
	// NamePL is the Polish name of the holiday.
	NamePL string
	// Date is the day on which the holiday is observed, at midnight UTC.
	Date time.Time
}

// holidayType distinguishes holidays that fall on the same calendar date
// every year from those that are computed relative to Easter Sunday.
type holidayType int

const (
	fixed holidayType = iota
	movable
)

// holidayRule is the static definition of a holiday, from which the concrete
// Holiday for a given year is derived.
type holidayRule struct {
	name   string
	namePL string
	typ    holidayType

	// month and day are used for fixed holidays.
	month time.Month
	day   int

	// afterEaster is the offset in days from Easter Sunday for movable holidays.
	afterEaster int
}
