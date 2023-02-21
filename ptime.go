package ptime

import (
	"fmt"
	"time"

	"github.com/blackchip-org/ptime/locale"
)

type P struct {
	Locale *locale.Locale
	Parser *Parser
}

func ForLocale(loc *locale.Locale) *P {
	return &P{
		Locale: loc,
		Parser: NewParser(loc),
	}
}

func ForLocaleName(name string) (*P, error) {
	loc, ok := locale.Lookup(name)
	if !ok {
		return nil, fmt.Errorf("unknown locale: %v", name)
	}
	return ForLocale(loc), nil
}

func (p *P) Parse(text string) (Parsed, error) {
	return p.Parser.Parse(text)
}

func (p *P) ParseDate(text string) (Parsed, error) {
	return p.Parser.ParseDate(text)
}

func (p *P) ParseTime(text string) (Parsed, error) {
	return p.Parser.ParseTime(text)
}

func (p *P) Time(parsed Parsed, now time.Time) (time.Time, error) {
	return Time(p.Locale, parsed, now)
}

func (p *P) Format(layout string, t time.Time) string {
	return Format(p.Locale, layout, t)
}
