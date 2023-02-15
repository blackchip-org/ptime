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
	OrdDay      string `json:",omitempty"`
	Hours       string `json:",omitempty"`
	Minutes     string `json:",omitempty"`
	Seconds     string `json:",omitempty"`
	FracSeconds string `json:",omitempty"`
	Period      string `json:",omitempty"`
	Zone        string `json:",omitempty"`
	Offset      string `json:",omitempty"`
	DateSep     string `json:",omitempty"`
	TimeSep     string `json:",omitempty"`
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
)

func (s state) String() string {
	switch s {
	case ParsingDate:
		return "ParsingDate"
	case ParsingTime:
		return "ParsingTime"
	case ParsingZone:
		return "ParsingZone"
	}
	return "Unknown"
}

type Parser struct {
	locale.Locale
	tokens []Token
	tok    Token
	idx    int
	parsed Parsed
	Trace  bool
	state  state
}

func NewParser(l locale.Locale) *Parser {
	return &Parser{Locale: l}
}

func (p *Parser) parse(text string) (Parsed, error) {
	p.tokens = Scan(text)
	p.idx = -1
	p.parsed = Parsed{}
	p.next()

	for p.tok.Type != End {
		var err error
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
		p.next()
	}
	return p.parsed, nil
}

func (p *Parser) Parse(text string) (Parsed, error) {
	p.changeState(ParsingDate)
	return p.parse(text)
}

func (p *Parser) ParseDate(text string) (Parsed, error) {
	p.changeState(ParsingDate)
	return p.parse(text)
}

func (p *Parser) ParseTime(text string) (Parsed, error) {
	p.changeState(ParsingTime)
	return p.parse(text)
}

func (p *Parser) parseText() error {
	if p.state == ParsingDate {
		if p.parsed.Weekday == "" {
			if day, ok := p.DayNames[p.tok.Val]; ok {
				p.parsed.Weekday = string(day)
				return nil
			}
		}
		if p.parsed.Month == "" {
			if mon, ok := p.MonthNames[p.tok.Val]; ok {
				p.parsed.Month = string(mon)
				return nil
			}
		}
	}
	if p.state == ParsingTime {
		if strings.Contains(p.TimeSep, p.tok.Val) {
			return nil
		}
		if _, ok := p.ZoneNames[p.tok.Val]; ok {
			p.changeState(ParsingZone)
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
	if p.state == ParsingDate {
		la := p.lookahead(1)
		if la.Type == Indicator && strings.Contains(p.TimeSep, la.Val) {
			p.changeState(ParsingTime)

		} else {
			return p.parseNumberDate()
		}
	}
	if p.state == ParsingTime {
		return p.parseNumberTime()
	}
	return nil
}

func (p *Parser) parseNumberDate() error {
	sep := p.parsed.DateSep
	if sep == "" {
		la := p.lookahead(1)
		if la.Type == Indicator {
			if strings.Contains(p.DateSep, la.Val) {
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
	if sep == "" {
		la := p.lookahead(1)
		p.trace("lookahead = '%v'", la.Val)
		if strings.Contains(p.TimeSep, la.Val) {
			sep = la.Val
		} else {
			sep = " "
		}
		p.trace("TimeSep = '%v'", sep)
		p.parsed.TimeSep = sep
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
		if p.parsed.Period == "" {
			period, ok := p.PeriodNames[p.tok.Val]
			if ok {
				p.trace("is period")
				p.parsed.Period = string(period)
				p.changeState(ParsingZone)
				return nil
			}
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
	if p.state != newState {
		p.trace("state: %v -> %v", p.state, newState)
	}
	p.state = newState
}

func (p *Parser) parseDate() error {
	delim := p.parsed.DateSep
	if delim == "-" {
		return p.parseYearDayMonth()
	}
	if delim != "" {
		la := p.lookahead(1)
		_, laIsMonth := p.MonthNames[la.Val]
		if la.Type == Text && laIsMonth {
			p.trace("lookahead(1) is '%v'", la.Val)
			return p.parseDayMonthYear()
		}
		if p.MonthDayOrder {
			return p.parseMonthDayYear()
		}
		return p.parseDayMonthYear()
	}
	return p.err("fixme")
}

func (p *Parser) parseYearDayMonth() error {
	p.trace("is YDM")
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

func (p *Parser) parseDayMonthYear() error {
	p.trace("is DMY")
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
	p.trace("is MDY")
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

func (p *Parser) parseTime() error {
	if p.parsed.Hours == "" {
		return p.parseHours()
	}
	if p.parsed.Minutes == "" {
		return p.parseMinutes()
	}
	if p.parsed.Seconds == "" {
		return p.parseSeconds()
	}
	return p.err("pass parseTime")
}

func (p *Parser) parseHours() error {
	p.trace("is hours")
	p.parsed.Hours = p.tok.Val
	h, err := strconv.Atoi(p.tok.Val)
	if err != nil {
		return p.err("invalid hours: %v", p.tok.Val)
	}
	if h < 0 || h >= 24 {
		return p.err("invalid hours: %v", p.tok.Val)
	}
	return nil
}

func (p *Parser) parseMinutes() error {
	p.trace("is minutes")
	p.parsed.Minutes = p.tok.Val
	m, err := strconv.Atoi(p.tok.Val)
	if err != nil {
		return p.err("invalid minutes: %v", p.tok.Val)
	}
	if m < 0 || m >= 60 {
		return p.err("invalid minutes: %v", p.tok.Val)
	}
	return nil
}

func (p *Parser) parseSeconds() error {
	p.trace("is seconds")
	p.parsed.Seconds = p.tok.Val
	s, err := strconv.Atoi(p.tok.Val)
	if err != nil {
		return p.err("invalid seconds: %v", p.tok.Val)
	}
	if s < 0 || s >= 60 {
		return p.err("invalid seconds: %v", p.tok.Val)
	}
	la := p.lookahead(1)
	if la.Type == Indicator && la.Val == p.DecimalSep {
		p.trace("has fractions")
		p.next()
		p.next()
		p.parsed.FracSeconds = p.tok.Val
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