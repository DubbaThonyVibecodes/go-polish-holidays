package polishholidays

import "time"

// holidayRules defines every public holiday in Poland.
var holidayRules = []holidayRule{
	{
		name:   "New Year",
		namePL: "Nowy Rok",
		typ:    fixed,
		month:  time.January,
		day:    1,
	},
	{
		name:   "Three Kings' Day",
		namePL: "Święto Trzech Króli",
		typ:    fixed,
		month:  time.January,
		day:    6,
	},
	{
		name:   "Labour Day",
		namePL: "Święto Pracy",
		typ:    fixed,
		month:  time.May,
		day:    1,
	},
	{
		name:   "3 May Constitution Day",
		namePL: "Narodowe Święto Konstytucji Trzeciego Maja",
		typ:    fixed,
		month:  time.May,
		day:    3,
	},
	{
		name:   "Assumption of Mary",
		namePL: "Wniebowzięcie Najświętszej Maryi Panny",
		typ:    fixed,
		month:  time.August,
		day:    15,
	},
	{
		name:   "All Saints' Day",
		namePL: "Wszystkich Świętych",
		typ:    fixed,
		month:  time.November,
		day:    1,
	},
	{
		name:   "National Independence Day",
		namePL: "Narodowe Święto Niepodległości",
		typ:    fixed,
		month:  time.November,
		day:    11,
	},
	{
		name:   "Christmas",
		namePL: "Boże Narodzenie",
		typ:    fixed,
		month:  time.December,
		day:    25,
	},
	{
		name:   "Second Day of Christmas",
		namePL: "Boże Narodzenie - drugi dzień",
		typ:    fixed,
		month:  time.December,
		day:    26,
	},
	{
		name:        "Easter Sunday",
		namePL:      "Niedziela Wielkanocna",
		typ:         movable,
		afterEaster: 0,
	},
	{
		name:        "Easter Monday",
		namePL:      "Poniedziałek Wielkanocny",
		typ:         movable,
		afterEaster: 1,
	},
	{
		name:        "Green Week",
		namePL:      "Zielone Świątki",
		typ:         movable,
		afterEaster: 49,
	},
	{
		name:        "Feast of Corpus Christi",
		namePL:      "Boże Ciało",
		typ:         movable,
		afterEaster: 60,
	},
}
