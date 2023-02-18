package locale

import (
	"fmt"
	"strings"
)

type String2D [][]string

func (s String2D) Main(index int) string {
	if len(s[index]) == 0 {
		return ""
	}
	return s[index][0]
}

func (s String2D) Alt(index int) string {
	switch len(s[index]) {
	case 0:
		return ""
	case 1:
		return s[index][0]
	default:
		return s[index][1]
	}
}

const (
	AM = iota
	PM
	Noon
	Midnight
)

type Def struct {
	MonthDayOrder     bool
	MonthNamesWide    []string
	MonthNamesAbbr    []string
	DayNamesWide      []string
	DayNamesAbbr      []string
	PeriodNamesAbbr   String2D
	PeriodNamesNarrow String2D
	ZoneNamesShort    map[string]string
	DateSep           []string
	TimeSep           []string
	HourSep           []string
	DecimalSep        string
	DateTimeSep       []string
	UTCFlags          []string
}

type Locale struct {
	Def
	MonthNum  map[string]int
	DayNum    map[string]int
	PeriodNum map[string]int
	Offsets   map[string]string
}

func New(def Def) (*Locale, error) {
	l := &Locale{
		Def:       def,
		MonthNum:  make(map[string]int),
		DayNum:    make(map[string]int),
		PeriodNum: make(map[string]int),
		Offsets:   make(map[string]string),
	}

	if len(def.MonthNamesAbbr) != 12 {
		return nil, fmt.Errorf("invalid number of month names (abbreviated)")
	}
	if len(def.MonthNamesWide) != 12 {
		return nil, fmt.Errorf("invalid number of month names (wide)")
	}
	for i := 0; i < 12; i++ {
		l.MonthNum[def.MonthNamesAbbr[i]] = i + 1
		l.MonthNum[def.MonthNamesWide[i]] = i + 1
		l.MonthNum[strings.ToLower(def.MonthNamesAbbr[i])] = i + 1
		l.MonthNum[strings.ToLower(def.MonthNamesWide[i])] = i + 1
	}
	// Probably need something better than this
	for k, v := range l.MonthNum {
		if strings.HasSuffix(k, ".") {
			l.MonthNum[k[:len(k)-1]] = v
		}
	}

	if len(def.DayNamesAbbr) != 7 {
		return nil, fmt.Errorf("invalid number of day names (abbreviated)")
	}
	if len(def.DayNamesWide) != 7 {
		return nil, fmt.Errorf("invalid number of day names (wide)")
	}
	for i := 0; i < 7; i++ {
		l.DayNum[def.DayNamesAbbr[i]] = i
		l.DayNum[def.DayNamesWide[i]] = i
		l.DayNum[strings.ToUpper(def.DayNamesAbbr[i])] = i
		l.DayNum[strings.ToUpper(def.DayNamesWide[i])] = i
	}

	for i, names := range def.PeriodNamesAbbr {
		for _, n := range names {
			l.PeriodNum[n] = i
			l.PeriodNum[strings.ToLower(n)] = i
		}
	}
	for i, names := range def.PeriodNamesNarrow {
		for _, n := range names {
			l.PeriodNum[n] = i
			l.PeriodNum[strings.ToLower(n)] = i
		}
	}

	for z, offset := range def.ZoneNamesShort {
		l.Offsets[z] = offset
		l.Offsets[strings.ToLower(z)] = offset
	}

	return l, nil
}

func MustNew(def Def) *Locale {
	l, err := New(def)
	if err != nil {
		panic(err)
	}
	return l
}

var table = map[string]*Locale{
	"en-US": EnUS,
	"fr-FR": FrFR,
}

func Lookup(name string) (*Locale, bool) {
	l, ok := table[name]
	return l, ok
}

func MustLookup(name string) *Locale {
	l, ok := Lookup(name)
	if !ok {
		panic(fmt.Sprintf("locale '%v' not found", name))
	}
	return l
}
