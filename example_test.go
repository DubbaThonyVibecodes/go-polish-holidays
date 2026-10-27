package polishholidays_test

import (
	"fmt"
	"time"

	polishholidays "github.com/DubbaThonyVibecodes/go-polish-holidays"
)

func ExampleIsHoliday() {
	labourDay := time.Date(2018, time.May, 1, 0, 0, 0, 0, time.UTC)
	fmt.Println(polishholidays.IsHoliday(labourDay))

	easter, _ := polishholidays.ParseDate("2020-04-12")
	fmt.Println(polishholidays.IsHoliday(easter))

	ordinary, _ := polishholidays.ParseDate("2020-12-03")
	fmt.Println(polishholidays.IsHoliday(ordinary))
	// Output:
	// true
	// true
	// false
}

func ExampleGetHolidayOnDate() {
	date, _ := polishholidays.ParseDate("2016-05-26")
	if holiday, ok := polishholidays.GetHolidayOnDate(date); ok {
		fmt.Println(holiday.Name)
	}

	date, _ = polishholidays.ParseDate("1978-10-02")
	_, ok := polishholidays.GetHolidayOnDate(date)
	fmt.Println(ok)
	// Output:
	// Feast of Corpus Christi
	// false
}

func ExampleGetHolidaysInYear() {
	holidays, err := polishholidays.GetHolidaysInYear(2020)
	if err != nil {
		panic(err)
	}
	fmt.Println(len(holidays))
	fmt.Println(holidays[0].Name)
	fmt.Println(holidays[0].NamePL)
	fmt.Println(holidays[0].Date.Format("2006-01-02"))
	// Output:
	// 13
	// New Year
	// Nowy Rok
	// 2020-01-01
}

func ExampleIsWorkingDay() {
	date, _ := polishholidays.ParseDate("2020-12-07") // Monday
	fmt.Println(polishholidays.IsWorkingDay(date))

	date, _ = polishholidays.ParseDate("2020-12-25") // Christmas
	fmt.Println(polishholidays.IsWorkingDay(date))
	// Output:
	// true
	// false
}
