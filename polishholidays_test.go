package polishholidays

import (
	"encoding/json"
	"errors"
	"os"
	"testing"
	"time"
)

func mustParse(t *testing.T, s string) time.Time {
	t.Helper()
	date, err := ParseDate(s)
	if err != nil {
		t.Fatalf("ParseDate(%q) returned unexpected error: %v", s, err)
	}
	return date
}

func TestParseDate(t *testing.T) {
	t.Run("accepts valid date strings", func(t *testing.T) {
		for _, s := range validDateStrings {
			if _, err := ParseDate(s); err != nil {
				t.Errorf("ParseDate(%q) = error %v, want nil", s, err)
			}
		}
	})

	t.Run("otherwise returns ErrInvalidDate", func(t *testing.T) {
		for _, s := range invalidDateStrings {
			_, err := ParseDate(s)
			if !errors.Is(err, ErrInvalidDate) {
				t.Errorf("ParseDate(%q) = error %v, want ErrInvalidDate", s, err)
			}
		}
	})

	t.Run("interprets date-only strings as midnight UTC", func(t *testing.T) {
		got := mustParse(t, "2020-03-20")
		want := time.Date(2020, time.March, 20, 0, 0, 0, 0, time.UTC)
		if !got.Equal(want) {
			t.Errorf("ParseDate(\"2020-03-20\") = %v, want %v", got, want)
		}
	})

	t.Run("preserves zone offsets", func(t *testing.T) {
		got := mustParse(t, "2020-12-11T02:06:04+01:00")
		want := time.Date(2020, time.December, 11, 1, 6, 4, 0, time.UTC)
		if !got.Equal(want) {
			t.Errorf("ParseDate(RFC3339) = %v, want %v", got, want)
		}
	})
}

func TestGetHolidayOnDate(t *testing.T) {
	t.Run("returns information about a holiday if there is a holiday on the provided date", func(t *testing.T) {
		for _, tc := range existingHolidays {
			holiday, ok := GetHolidayOnDate(mustParse(t, tc.date))
			if !ok {
				t.Errorf("GetHolidayOnDate(%s): no holiday found, want %q", tc.date, tc.name)
				continue
			}
			if holiday.Name != tc.name {
				t.Errorf("GetHolidayOnDate(%s).Name = %q, want %q", tc.date, holiday.Name, tc.name)
			}
			if holiday.NamePL == "" {
				t.Errorf("GetHolidayOnDate(%s).NamePL is empty", tc.date)
			}
			if holiday.Date.IsZero() {
				t.Errorf("GetHolidayOnDate(%s).Date is zero", tc.date)
			}
		}
	})

	t.Run("returns false if there is no holiday on the provided date", func(t *testing.T) {
		for _, date := range nonHolidays {
			if holiday, ok := GetHolidayOnDate(mustParse(t, date)); ok {
				t.Errorf("GetHolidayOnDate(%s) = %+v, want no holiday", date, holiday)
			}
		}
	})

	t.Run("returns false for dates before year 1", func(t *testing.T) {
		date := time.Date(0, time.January, 1, 0, 0, 0, 0, time.UTC)
		if _, ok := GetHolidayOnDate(date); ok {
			t.Errorf("GetHolidayOnDate(year 0) found a holiday, want none")
		}
	})

	t.Run("compares calendar days in the input's own location", func(t *testing.T) {
		warsaw, err := time.LoadLocation("Europe/Warsaw")
		if err != nil {
			t.Skipf("tzdata unavailable: %v", err)
		}
		losAngeles, err := time.LoadLocation("America/Los_Angeles")
		if err != nil {
			t.Skipf("tzdata unavailable: %v", err)
		}

		cases := []struct {
			date time.Time
			want bool
		}{
			{time.Date(2020, time.April, 12, 23, 59, 59, 0, warsaw), true},      // Easter Sunday, late evening in Poland
			{time.Date(2020, time.April, 12, 0, 0, 1, 0, losAngeles), true},     // Easter Sunday, just after midnight in LA
			{time.Date(2020, time.April, 13, 0, 0, 0, 0, warsaw), true},         // Easter Monday
			{time.Date(2020, time.April, 11, 23, 59, 59, 0, losAngeles), false}, // Saturday before Easter in LA
			{time.Date(2020, time.April, 14, 0, 0, 0, 0, warsaw), false},
		}
		for _, tc := range cases {
			if got := IsHoliday(tc.date); got != tc.want {
				t.Errorf("IsHoliday(%v) = %v, want %v", tc.date, got, tc.want)
			}
		}
	})
}

func TestIsHoliday(t *testing.T) {
	t.Run("returns true if there is a holiday on the provided date", func(t *testing.T) {
		for _, tc := range existingHolidays {
			if !IsHoliday(mustParse(t, tc.date)) {
				t.Errorf("IsHoliday(%s) = false, want true", tc.date)
			}
		}
	})

	t.Run("returns false if there is no holiday on the provided date", func(t *testing.T) {
		for _, date := range nonHolidays {
			if IsHoliday(mustParse(t, date)) {
				t.Errorf("IsHoliday(%s) = true, want false", date)
			}
		}
	})
}

func TestGetHolidaysInYear(t *testing.T) {
	t.Run("accepts positive years", func(t *testing.T) {
		for _, year := range validYears {
			if _, err := GetHolidaysInYear(year); err != nil {
				t.Errorf("GetHolidaysInYear(%d) = error %v, want nil", year, err)
			}
		}
	})

	t.Run("otherwise returns ErrInvalidYear", func(t *testing.T) {
		for _, year := range invalidYears {
			holidays, err := GetHolidaysInYear(year)
			if !errors.Is(err, ErrInvalidYear) {
				t.Errorf("GetHolidaysInYear(%d) = error %v, want ErrInvalidYear", year, err)
			}
			if holidays != nil {
				t.Errorf("GetHolidaysInYear(%d) returned %d holidays alongside an error", year, len(holidays))
			}
		}
	})

	t.Run("returns correct dates of public holidays for the given year", func(t *testing.T) {
		for _, tc := range holidaysInYears {
			holidays, err := GetHolidaysInYear(tc.year)
			if err != nil {
				t.Fatalf("GetHolidaysInYear(%d): %v", tc.year, err)
			}
			found := false
			for _, holiday := range holidays {
				if holiday.Name != tc.holidayName {
					continue
				}
				found = true
				if got := holiday.Date.Format("2006-01-02"); got != tc.holidayDate {
					t.Errorf("%d %s = %s, want %s", tc.year, tc.holidayName, got, tc.holidayDate)
				}
			}
			if !found {
				t.Errorf("GetHolidaysInYear(%d) does not contain %q", tc.year, tc.holidayName)
			}
		}
	})

	t.Run("returns names of public holidays in English and Polish", func(t *testing.T) {
		for _, year := range yearsInPastAndFuture {
			holidays, err := GetHolidaysInYear(year)
			if err != nil {
				t.Fatalf("GetHolidaysInYear(%d): %v", year, err)
			}
			for _, holiday := range holidays {
				if holiday.Name == "" || holiday.NamePL == "" {
					t.Errorf("%d: holiday %+v has an empty name", year, holiday)
				}
				if holiday.Name == holiday.NamePL {
					t.Errorf("%d: holiday %q has identical English and Polish names", year, holiday.Name)
				}
			}
		}
	})

	t.Run("returns valid dates for years in the past and future", func(t *testing.T) {
		for _, year := range yearsInPastAndFuture {
			holidays, err := GetHolidaysInYear(year)
			if err != nil {
				t.Fatalf("GetHolidaysInYear(%d): %v", year, err)
			}
			for _, holiday := range holidays {
				if holiday.Date.IsZero() {
					t.Errorf("%d: %q has a zero date", year, holiday.Name)
				}
				if holiday.Date.Year() != year {
					t.Errorf("%d: %q falls in year %d", year, holiday.Name, holiday.Date.Year())
				}
				if holiday.Date.Location() != time.UTC {
					t.Errorf("%d: %q is not in UTC", year, holiday.Name)
				}
				if h, m, s := holiday.Date.Clock(); h != 0 || m != 0 || s != 0 {
					t.Errorf("%d: %q is not at midnight: %v", year, holiday.Name, holiday.Date)
				}
			}
		}
	})

	t.Run("returns all 13 holidays sorted by date", func(t *testing.T) {
		holidays, err := GetHolidaysInYear(2020)
		if err != nil {
			t.Fatal(err)
		}
		if len(holidays) != 13 {
			t.Fatalf("len = %d, want 13", len(holidays))
		}
		if holidays[0].Name != "New Year" || holidays[0].NamePL != "Nowy Rok" {
			t.Errorf("first holiday = %+v, want New Year / Nowy Rok", holidays[0])
		}
		if last := holidays[len(holidays)-1]; last.Name != "Second Day of Christmas" {
			t.Errorf("last holiday = %+v, want Second Day of Christmas", last)
		}
		for i := 1; i < len(holidays); i++ {
			if !holidays[i-1].Date.Before(holidays[i].Date) {
				t.Errorf("holidays not sorted: %v before %v", holidays[i-1].Date, holidays[i].Date)
			}
		}
	})

	t.Run("matches the original JavaScript library for 1900-2100", func(t *testing.T) {
		raw, err := os.ReadFile("testdata/holidays_1900_2100.json")
		if err != nil {
			t.Fatal(err)
		}
		var golden map[string][]struct {
			Name   string `json:"name"`
			NamePL string `json:"namePL"`
			Date   string `json:"date"`
		}
		if err := json.Unmarshal(raw, &golden); err != nil {
			t.Fatal(err)
		}
		if len(golden) != 201 {
			t.Fatalf("golden file has %d years, want 201", len(golden))
		}
		for year := 1900; year <= 2100; year++ {
			want := golden[time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC).Format("2006")]
			got, err := GetHolidaysInYear(year)
			if err != nil {
				t.Fatalf("GetHolidaysInYear(%d): %v", year, err)
			}
			if len(got) != len(want) {
				t.Errorf("%d: got %d holidays, want %d", year, len(got), len(want))
				continue
			}
			for i := range want {
				if got[i].Name != want[i].Name || got[i].NamePL != want[i].NamePL || got[i].Date.Format("2006-01-02") != want[i].Date {
					t.Errorf("%d[%d]: got %s/%s/%s, want %s/%s/%s", year, i,
						got[i].Name, got[i].NamePL, got[i].Date.Format("2006-01-02"),
						want[i].Name, want[i].NamePL, want[i].Date)
				}
			}
		}
	})
}

func TestWeekendAndWorkingDay(t *testing.T) {
	cases := []struct {
		date             string
		weekend, holiday bool
		weekendOrHoliday bool
		workingDay       bool
	}{
		{"2020-12-05", true, false, true, false},  // Saturday
		{"2020-12-06", true, false, true, false},  // Sunday
		{"2020-12-07", false, false, false, true}, // Monday
		{"2020-12-25", false, true, true, false},  // Christmas on a Friday
		{"2020-04-12", true, true, true, false},   // Easter Sunday
		{"2020-04-13", false, true, true, false},  // Easter Monday
		{"2021-05-01", true, true, true, false},   // Labour Day on a Saturday
	}
	for _, tc := range cases {
		date := mustParse(t, tc.date)
		if got := IsWeekend(date); got != tc.weekend {
			t.Errorf("IsWeekend(%s) = %v, want %v", tc.date, got, tc.weekend)
		}
		if got := IsHoliday(date); got != tc.holiday {
			t.Errorf("IsHoliday(%s) = %v, want %v", tc.date, got, tc.holiday)
		}
		if got := IsWeekendOrHoliday(date); got != tc.weekendOrHoliday {
			t.Errorf("IsWeekendOrHoliday(%s) = %v, want %v", tc.date, got, tc.weekendOrHoliday)
		}
		if got := IsWorkingDay(date); got != tc.workingDay {
			t.Errorf("IsWorkingDay(%s) = %v, want %v", tc.date, got, tc.workingDay)
		}
	}
}
