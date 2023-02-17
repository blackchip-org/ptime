package ptime

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/blackchip-org/ptime/locale"
)

type Parsed struct {
	Weekday     string `json:",omitempty"`
	Year        string `json:",omitempty"`
	Month       string `json:",omitempty"`
	Day         string `json:",omitempty"`
	Hour        string `json:",omitempty"`
	Minute      string `json:",omitempty"`
	Second      string `json:",omitempty"`
	FracSecond  string `json:",omitempty"`
	Period      string `json:",omitempty"`
	Zone        string `json:",omitempty"`
	Offset      string `json:",omitempty"`
	DateSep     string `json:",omitempty"`
	TimeSep     string `json:",omitempty"`
	DateTimeSep string `json:",omitempty"`
	HourSep     string `json:",omitempty"`
}

func (p Parsed) String() string {
	text, err := json.Marshal(p)
	if err != nil {
		panic(err)
	}
	return string(text)
}

type state int

const (
	Unknown state = iota
	ParsingDate
	ParsingTime
	ParsingZone
	Done
)

func (s state) String() string {
	switch s {
	case ParsingDate:
		return "ParsingDate"
	case ParsingTime:
		return "ParsingTime"
	case ParsingZone:
		return "ParsingZone"
	case Done:
		return "Done"
	}
	return "Unknown"
}

type dateOrder int

const (
	UnknownOrder dateOrder = iota
	DayMonthYearOrder
	MonthDayYearOrder
	YearMonthDayOrder
	YearDayOrder
)

func (d dateOrder) String() string {
	switch d {
	case DayMonthYearOrder:
		return "day-month-year"
	case MonthDayYearOrder:
		return "month-day-year"
	case YearMonthDayOrder:
		return "year-month-day"
	case YearDayOrder:
		return "year-day"
	}
	return "unknown"
}

type Parser struct {
	locale.Locale
	tokens    []Token
	tok       Token
	idx       int
	parsed    Parsed
	Trace     bool
	state     state
	dateOrder dateOrder
	parseOne  bool
}

func NewParser(l locale.Locale) *Parser {
	return &Parser{Locale: l}
}

func (p *Parser) parse(text string) (Parsed, error) {
	p.trace("state: %v", p.state)
	text = strings.ToLower(text)
	for i := 0; i < len(p.Replacements); i += 2 {
		from := p.Replacements[i]
		to := p.Replacements[i+1]
		text = strings.ReplaceAll(text, from, to)
	}
	p.tokens = Scan(text)

	if len(p.tokens) == 0 {
		return Parsed{}, nil
	}

	p.idx = -1
	p.tok = p.tokens[0]
	p.parsed = Parsed{}
	p.dateOrder = UnknownOrder

	for p.tok.Type != End {
		var err error
		p.trace("top")
		p.next()
		switch p.tok.Type {
		case Text:
			err = p.parseText()
		case Number:
			err = p.parseNumber()
		case Indicator:
			err = p.parseIndicator()
		}
		if err != nil {
			return p.parsed, err
		}
	}
	return p.parsed, nil
}

func (p *Parser) Parse(text string) (Parsed, error) {
	p.state = Unknown
	p.parseOne = false
	return p.parse(text)
}

func (p *Parser) ParseDate(text string) (Parsed, error) {
	p.state = ParsingDate
	p.parseOne = true
	return p.parse(text)
}

func (p *Parser) ParseTime(text string) (Parsed, error) {
	p.state = ParsingTime
	p.parseOne = true
	return p.parse(text)
}

func (p *Parser) parseText() error {
	if p.state == Unknown {
		p.state = ParsingDate
	}
	if p.state == ParsingDate {
		if p.parsed.Weekday == "" {
			if day, ok := p.DayNames[p.tok.Val]; ok {
				p.trace("is weekday")
				p.parsed.Weekday = string(day)
				return nil
			}
		}
		if p.parsed.Month == "" {
			if mon, ok := p.MonthNames[p.tok.Val]; ok {
				p.trace("is month")
				p.parsed.Month = string(mon)
				return nil
			}
		}
		if is(p.tok.Val, p.DateTimeSep) {
			p.changeState(ParsingTime)
			p.parsed.DateTimeSep = p.tok.Val
			return nil
		}
	}
	if p.state == ParsingTime {
		if is(p.tok.Val, p.TimeSep) {
			return nil
		}
		if p.parsed.Period == "" {
			period, ok := p.PeriodNames[p.tok.Val]
			if ok {
				p.trace("is period")
				p.parsed.Period = string(period)
				p.changeState(ParsingZone)
				return nil
			}
		}
		if _, ok := p.ZoneNames[p.tok.Val]; ok {
			p.changeState(ParsingZone)
		}
		if is(p.tok.Val, p.UTCFlags) {
			p.parsed.Zone = p.tok.Val
			p.parsed.Offset = "+0000"
			p.changeState(Done)
			return nil
		}
	}
	if p.state == ParsingZone {
		p.parsed.Zone = p.tok.Val
		offset, ok := p.ZoneNames[p.tok.Val]
		if !ok {
			return p.err("unknown zone: %v", p.tok.Val)
		}
		if p.parsed.Offset != "" && p.parsed.Offset != offset {
			return p.err("time zone '%v' does not match given offset '%v'", p.tok.Val, p.parsed.Offset)
		}
		p.parsed.Offset = offset
		return nil
	}

	return p.err("unexpected text: %v", p.tok.Val)
}

func (p *Parser) parseNumber() error {
	if p.state == Unknown {
		la := p.lookahead(1)
		if la.Type == Indicator && (is(la.Val, p.TimeSep) || is(la.Val, p.HourSep)) {
			p.changeState(ParsingTime)
		} else {
			p.changeState(ParsingDate)
		}
	}
	if p.state == ParsingDate {
		la := p.lookahead(1)
		if la.Type == Indicator && is(la.Val, p.TimeSep) {
			p.changeState(ParsingTime)

		} else {
			return p.parseNumberDate()
		}
	}
	if p.state == ParsingTime {
		return p.parseNumberTime()
	}
	if p.state == ParsingZone {
		p.changeState(Done)
		return p.parseYear4()
	}
	return p.err("extra number: %v", p.tok.Val)
}

func (p *Parser) parseNumberDate() error {
	sep := p.parsed.DateSep
	if sep == "" {
		la := p.lookahead(1)
		if la.Type == Indicator {
			if is(la.Val, p.DateSep) {
				sep = la.Val
			}
		} else {
			sep = " "
		}
		p.trace("DateSep = '%v'", sep)
		p.parsed.DateSep = sep
	}
	return p.parseDate()
}

func (p *Parser) parseNumberTime() error {
	sep := p.parsed.TimeSep
	if sep == "" && p.parsed.HourSep == "" {
		la := p.lookahead(1)
		if la.Val != "" {
			if is(la.Val, p.TimeSep) {
				sep = la.Val
			} else {
				if is(la.Val, p.HourSep) {
					sep = ""
				}
			}
			p.trace("TimeSep = '%v'", sep)
			p.parsed.TimeSep = sep
		}
	}
	return p.parseTime()
}

func (p *Parser) parseIndicator() error {
	if p.state == ParsingDate && p.tok.Val == p.parsed.DateSep {
		p.next()
		return p.parseDate()
	}
	if p.state == ParsingTime {
		if p.tok.Val == p.parsed.TimeSep {
			p.next()
			return p.parseTime()
		}
		if p.tok.Val == "-" || p.tok.Val == "+" {
			p.changeState(ParsingZone)
		}
	}
	if p.state == ParsingZone {
		if p.tok.Val == "-" || p.tok.Val == "+" {
			return p.parseOffset()
		}
	}
	p.trace("discarding")
	return nil
}

func (p *Parser) changeState(newState state) {
	if p.parseOne {
		if newState != ParsingZone {
			newState = Done
		}
	}
	if p.state != newState {
		p.trace("state: %v -> %v", p.state, newState)
	}
	p.state = newState
}

func (p *Parser) parseDate() error {
	delim := p.parsed.DateSep
	if p.dateOrder == UnknownOrder {
		la1 := p.lookahead(1)
		_, la1IsMonth := p.MonthNames[la1.Val]
		la2 := p.lookahead(2)
		_, la2IsMonth := p.MonthNames[la2.Val]

		//p.trace("lookahead: %v", la.Val)
		switch {
		case delim == "-" && la2IsMonth:
			p.dateOrder = DayMonthYearOrder
		case delim == "-" && la2.Type == Number && len(la2.Val) == 3:
			p.dateOrder = YearDayOrder
		case delim == "-":
			p.dateOrder = YearMonthDayOrder
		case la1IsMonth:
			p.dateOrder = DayMonthYearOrder
		case p.MonthDayOrder:
			p.dateOrder = MonthDayYearOrder
		default:
			p.dateOrder = DayMonthYearOrder
		}
		p.trace("order: %v", p.dateOrder)
	}
	switch p.dateOrder {
	case YearDayOrder:
		return p.parseYearDay()
	case YearMonthDayOrder:
		return p.parseYearMonthDay()
	case DayMonthYearOrder:
		return p.parseDayMonthYear()
	case MonthDayYearOrder:
		return p.parseMonthDayYear()
	}
	return p.err("unexpected '%v' in date", p.tok.Val)
}

func (p *Parser) parseYearMonthDay() error {
	if p.parsed.Year == "" {
		return p.parseYear4()
	}
	if p.parsed.Month == "" {
		return p.parseMonth()
	}
	if p.parsed.Day == "" {
		return p.parseDay()
	}
	return p.err("pass parseYearDayMonth")
}

func (p *Parser) parseYearDay() error {
	if p.parsed.Year == "" {
		return p.parseYear4()
	}
	if p.parsed.Day == "" {
		return p.parseOrdinalDay()
	}
	return p.err("pass parseYearDayMonth")
}

func (p *Parser) parseDayMonthYear() error {
	if p.parsed.Day == "" {
		return p.parseDay()
	}
	if p.parsed.Month == "" {
		return p.parseMonth()
	}
	if p.parsed.Year == "" {
		return p.parseYear()
	}
	return p.err("pass parseDayMonth")
}

func (p *Parser) parseMonthDayYear() error {
	if p.parsed.Month == "" {
		return p.parseMonth()
	}
	if p.parsed.Day == "" {
		return p.parseDay()
	}
	if p.parsed.Year == "" {
		return p.parseYear()
	}
	return p.err("pass parseMonthDay")
}

func (p *Parser) parseYear() error {
	p.trace("is year")
	p.parsed.Year = p.tok.Val
	if len(p.parsed.Year) != 4 && len(p.parsed.Year) != 2 {
		return p.err("invalid year: %v", p.parsed.Year)
	}
	return nil

}

func (p *Parser) parseYear4() error {
	p.trace("is year4")
	p.parsed.Year = p.tok.Val
	if len(p.parsed.Year) != 4 {
		return p.err("invalid year: %v", p.parsed.Year)
	}
	return nil
}

func (p *Parser) parseMonth() error {
	p.trace("is month")
	p.parsed.Month = p.tok.Val
	if _, ok := p.MonthNames[p.tok.Val]; ok {
		return nil
	}
	m, err := strconv.Atoi(p.tok.Val)
	if err != nil {
		return p.err("invalid month: %v", p.tok.Val)
	}
	if m < 1 || m > 12 {
		return p.err("invalid month: %v", p.tok.Val)
	}
	return nil
}

func (p *Parser) parseDay() error {
	p.trace("is day")
	p.parsed.Day = p.tok.Val
	d, err := strconv.Atoi(p.tok.Val)
	if err != nil {
		return p.err("invalid day: %v", p.tok.Val)
	}
	if d < 1 || d > 31 {
		return p.err("invalid day: %v", p.tok.Val)
	}
	return nil
}

func (p *Parser) parseOrdinalDay() error {
	p.trace("is ordinal day")
	p.parsed.Day = p.tok.Val
	d, err := strconv.Atoi(p.tok.Val)
	if err != nil {
		return p.err("invalid day: %v", p.tok.Val)
	}
	if d < 1 || d > 365 {
		return p.err("invalid day: %v", p.tok.Val)
	}
	return nil
}

func (p *Parser) parseTime() error {
	if p.parsed.Hour == "" {
		return p.parseHour()
	}
	if p.parsed.Minute == "" {
		return p.parseMinute()
	}
	if p.parsed.Second == "" {
		return p.parseSecond()
	}
	p.changeState(Done)
	return p.parseYear4()
}

func (p *Parser) parseHour() error {
	p.trace("is hour")
	p.parsed.Hour = p.tok.Val
	h, err := strconv.Atoi(p.tok.Val)
	if err != nil {
		return p.err("invalid hour: %v", p.tok.Val)
	}
	if h < 0 || h >= 24 {
		return p.err("invalid hour: %v", p.tok.Val)
	}
	la := p.lookahead(1)
	if is(la.Val, p.HourSep) {
		p.trace("HourSep = '%v'", la.Val)
		p.parsed.HourSep = la.Val
		p.next()
	}
	return nil
}

func (p *Parser) parseMinute() error {
	p.trace("is minute")
	p.parsed.Minute = p.tok.Val
	m, err := strconv.Atoi(p.tok.Val)
	if err != nil {
		return p.err("invalid minute: %v", p.tok.Val)
	}
	if m < 0 || m >= 60 {
		return p.err("invalid minute: %v", p.tok.Val)
	}
	return nil
}

func (p *Parser) parseSecond() error {
	p.trace("is second")
	p.parsed.Second = p.tok.Val
	s, err := strconv.Atoi(p.tok.Val)
	if err != nil {
		return p.err("invalid second: %v", p.tok.Val)
	}
	if s < 0 || s >= 60 {
		return p.err("invalid second: %v", p.tok.Val)
	}
	la := p.lookahead(1)
	if la.Type == Indicator && la.Val == p.DecimalSep {
		p.trace("has fractions")
		p.next()
		p.next()
		p.parsed.FracSecond = p.tok.Val
	}
	return nil
}

func (p *Parser) parseOffset() error {
	p.trace("is offset")
	var parts []string
	if p.tok.Type == Indicator && (p.tok.Val == "+" || p.tok.Val == "-") {
		parts = append(parts, p.tok.Val)
		p.next()
	}
	if p.tok.Type != Number {
		return p.err("expecting offset but got '%v'", p.tok.Val)
	}
	if len(p.tok.Val) == 4 {
		parts = append(parts, p.tok.Val)
	}
	if len(p.tok.Val) == 2 {
		parts = append(parts, p.tok.Val)
		p.next()
		if p.tok.Type != Indicator || p.tok.Val != ":" {
			return p.err("expecting ':' in offset but got '%v'", p.tok.Val)
		}
		p.next()
		if p.tok.Type != Number {
			return p.err("expecting offset minutes but got '%v'", p.tok.Val)
		}
		parts = append(parts, p.tok.Val)
	}
	offset := strings.Join(parts, "")
	if p.parsed.Offset != "" && p.parsed.Offset != offset {
		return p.err("offset mismatch between '%v' and '%v'", offset, p.parsed.Offset)
	}
	p.parsed.Offset = offset
	return nil
}

func (p *Parser) lookahead(n int) Token {
	if n+p.idx >= len(p.tokens) {
		return Token{End, "", 0}
	}
	return p.tokens[n+p.idx]
}

func (p *Parser) next() {
	p.idx++
	if p.idx >= len(p.tokens) {
		p.trace("end")
		p.idx = len(p.tokens)
		p.tok = Token{End, "", 0}
		return
	}
	p.tok = p.tokens[p.idx]
	p.trace("next: %v", p.tok)
}

func (p *Parser) err(format string, a ...any) error {
	return fmt.Errorf(format, a...)
}

func (p *Parser) trace(format string, a ...any) {
	if p.Trace {
		fmt.Printf(format, a...)
		fmt.Println()
	}
}

func is(text string, domain []string) bool {
	for _, v := range domain {
		if text == v {
			return true
		}
	}
	return false
}
