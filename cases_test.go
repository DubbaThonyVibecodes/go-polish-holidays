package polishholidays

// Test case data ported from the original library's spec/holidays.cases.ts.

var invalidDateStrings = []string{
	"notavaliddate",
	"16",
	"",
	"2020-13-01",
	"2020-02-30",
	"1607648764000", // Date.now() as a string
}

var invalidYears = []int{0, -1, -500}

var validDateStrings = []string{
	"2020-03-20",
	"2020-12-11T02:06:04+01:00",
	"2020-12-11T02:06:04Z",
	"2020-12-11T02:06:04.123Z",
	"2020-12-11T02:06:04",
	"2020-12-11T02:06",
	"2020-12-11 02:06:04",
	"2019-07-21",
}

var validYears = []int{2020, 1950, 150, 3827, 15}

var holidaysInYears = []struct {
	year        int
	holidayName string
	holidayDate string
}{
	{2020, "Easter Sunday", "2020-04-12"},
	{2050, "Easter Sunday", "2050-04-10"},
	{2015, "Green Week", "2015-05-24"},
	{1950, "Feast of Corpus Christi", "1950-06-08"},
	{1994, "Labour Day", "1994-05-01"},
}

var yearsInPastAndFuture = []int{32, 299, 603, 1271, 1700, 1893, 1944, 2001, 2920, 3533}

var existingHolidays = []struct {
	date string
	name string
}{
	{"2006-01-06", "Three Kings' Day"},
	{"1950-05-28", "Green Week"},
	{"1990-12-25", "Christmas"},
}

var nonHolidays = []string{"1978-10-02", "2021-02-13", "2020-12-07", "1919-06-24"}

// easterDates holds known Easter Sunday dates for spot-checking the algorithm.
var easterDates = map[int]string{
	1583: "1583-04-10",
	1818: "1818-03-22", // earliest possible date
	1943: "1943-04-25", // latest possible date
	2000: "2000-04-23",
	2024: "2024-03-31",
	2025: "2025-04-20",
	2026: "2026-04-05",
	2038: "2038-04-25",
	2285: "2285-03-22",
}
