package locale

import (
	"fmt"
	"strconv"
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
	Offsets   map[string]int
}

func New(def Def) (*Locale, error) {
	l := &Locale{
		Def:       def,
		MonthNum:  make(map[string]int),
		DayNum:    make(map[string]int),
		PeriodNum: make(map[string]int),
		Offsets:   make(map[string]int),
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
		runes := []rune(offset)
		if len(runes) != 5 {
			return nil, fmt.Errorf("invalid offset: %v", offset)
		}

		var sign int
		switch runes[0] {
		case '+':
			sign = 1
		case '-':
			sign = -1
		default:
			return nil, fmt.Errorf("invalid offset: %v", offset)
		}

		hrs, err := strconv.Atoi(string(runes[1:3]))
		if err != nil {
			return nil, fmt.Errorf("invalid offset: %v", offset)
		}
		min, err := strconv.Atoi(string(runes[3:5]))
		if err != nil {
			return nil, fmt.Errorf("invalid offset: %v", offset)
		}
		offset := sign*hrs*3600 + min*60
		l.Offsets[z] = offset
		l.Offsets[strings.ToLower(z)] = offset
	}
	for _, flag := range l.UTCFlags {
		l.Offsets[flag] = 0
		l.Offsets[strings.ToLower(flag)] = 0
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
