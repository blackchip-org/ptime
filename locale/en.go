package locale

var EnMonthNames = map[string]Month{
	"january":   Jan,
	"jan":       Jan,
	"february":  Feb,
	"feb":       Feb,
	"march":     Mar,
	"mar":       Mar,
	"april":     Apr,
	"apr":       Apr,
	"may":       May,
	"june":      Jun,
	"jun":       Jun,
	"july":      Jul,
	"jul":       Jul,
	"august":    Aug,
	"aug":       Aug,
	"september": Sep,
	"sep":       Sep,
	"october":   Oct,
	"oct":       Oct,
	"november":  Nov,
	"nov":       Nov,
	"december":  Dec,
	"dec":       Dec,
}

var EnDayNames = map[string]Day{
	"sunday":    Sun,
	"sun":       Sun,
	"monday":    Mon,
	"mon":       Mon,
	"tuesday":   Tue,
	"tue":       Tue,
	"wednesday": Wed,
	"wed":       Wed,
	"thursday":  Thu,
	"thr":       Thu,
	"friday":    Fri,
	"fri":       Fri,
	"saturday":  Sat,
	"sat":       Sat,
}

var EnPeriodNames = map[string]Period{
	"am":       AM,
	"a.m.":     AM,
	"pm":       PM,
	"p.m.":     PM,
	"noon":     Noon,
	"midnight": Midnight,
}

var USZones = map[string]string{
	"est": "-0500",
	"cst": "-0600",
	"mst": "-0700",
	"pst": "-0800",
	"edt": "-0400",
	"cdt": "-0500",
	"mdt": "-0600",
	"pdt": "-0700",
}

var EnUS = Locale{
	MonthDayOrder: true,
	MonthNames:    EnMonthNames,
	DayNames:      EnDayNames,
	PeriodNames:   EnPeriodNames,
	ZoneNames:     USZones,
	DateSep:       "-/",
	TimeSep:       ":",
	DecimalSep:    ".",
}
