package locale

var FrMonthNames = map[string]Month{
	"janvier":   Jan,
	"janv":      Jan,
	"février":   Feb,
	"févr":      Feb,
	"mars":      Mar,
	"avril":     Apr,
	"avr":       Apr,
	"mai":       May,
	"juin":      Jun,
	"juillet":   Jul,
	"juil":      Jul,
	"août":      Aug,
	"septembre": Sep,
	"sept":      Sep,
	"octobre":   Oct,
	"oct":       Oct,
	"novembre":  Nov,
	"nov":       Nov,
	"décembre":  Dec,
	"déc":       Dec,
}

var FrDayNames = map[string]Day{
	"dimanche": Sun,
	"dim":      Sun,
	"lundi":    Mon,
	"lun":      Mon,
	"mardi":    Tue,
	"mar":      Tue,
	"mercredi": Wed,
	"mer":      Wed,
	"jeudi":    Thu,
	"jeu":      Thu,
	"vendredi": Fri,
	"ven":      Fri,
	"samedi":   Sat,
	"sam":      Sat,
}

var FrFR = Locale{
	MonthNames:  FrMonthNames,
	DayNames:    FrDayNames,
	DateSep:     []string{"-", "/"},
	TimeSep:     []string{":"},
	HourSep:     []string{"h"},
	DecimalSep:  ",",
	DateTimeSep: []string{"t"},
	UTCFlags:    []string{"utc", "z"},
}
