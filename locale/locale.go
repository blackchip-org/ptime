package locale

import "fmt"

type Month string

const (
	Jan Month = "jan"
	Feb Month = "feb"
	Mar Month = "mar"
	Apr Month = "apr"
	May Month = "may"
	Jun Month = "jun"
	Jul Month = "jul"
	Aug Month = "aug"
	Sep Month = "sep"
	Oct Month = "oct"
	Nov Month = "nov"
	Dec Month = "dec"
)

type Day string

const (
	Sun Day = "sun"
	Mon Day = "mon"
	Tue Day = "tue"
	Wed Day = "wed"
	Thu Day = "thu"
	Fri Day = "fri"
	Sat Day = "sat"
)

type Period string

const (
	AM       Period = "a"
	PM       Period = "p"
	Noon     Period = "n"
	Midnight Period = "m"
)

type Locale struct {
	MonthDayOrder bool
	MonthNames    map[string]Month
	DayNames      map[string]Day
	PeriodNames   map[string]Period
	ZoneNames     map[string]string
	DateSep       map[string]struct{}
	TimeSep       map[string]struct{}
	HourSep       map[string]struct{}
	DecimalSep    string
}

var HyphenSlashSep = map[string]struct{}{
	"-": {},
	"/": {},
}

var ColonSep = map[string]struct{}{
	":": {},
}

var table = map[string]Locale{
	"en-US": EnUS,
	"fr-FR": FrFR,
}

func Lookup(name string) (Locale, bool) {
	l, ok := table[name]
	return l, ok
}

func MustLookup(name string) Locale {
	l, ok := Lookup(name)
	if !ok {
		panic(fmt.Sprintf("locale '%v' not found", name))
	}
	return l
}
