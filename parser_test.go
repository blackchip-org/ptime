package ptime

import (
	"strings"
	"testing"

	"github.com/blackchip-org/ptime/locale"
)

func TestParserEnUS(t *testing.T) {
	tests := []struct {
		fn     string
		text   string
		parsed Parsed
	}{
		{"date", "2006-01-02", Parsed{Year: "2006", Month: "01", Day: "02", DateSep: "-"}},
		{"date", "2006-01", Parsed{Year: "2006", Month: "01", DateSep: "-"}},
		{"date", "1/2", Parsed{Month: "1", Day: "2", DateSep: "/"}},
		{"date", "1/2/2006", Parsed{Month: "1", Day: "2", Year: "2006", DateSep: "/"}},
		{"date", "1/2/06", Parsed{Month: "1", Day: "2", Year: "06", DateSep: "/"}},
		{"date", "Jan 2 2006", Parsed{Month: "jan", Day: "2", Year: "2006", DateSep: " "}},
		{"date", "Jan 2 06", Parsed{Month: "jan", Day: "2", Year: "06", DateSep: " "}},
		{"date", "Mon Jan 2 2006", Parsed{Weekday: "mon", Month: "jan", Day: "2", Year: "2006", DateSep: " "}},
		{"date", "Monday Jan 2 2006", Parsed{Weekday: "mon", Month: "jan", Day: "2", Year: "2006", DateSep: " "}},
		{"date", "Jan 2", Parsed{Month: "jan", Day: "2", DateSep: " "}},
		{"date", "Mon, Jan 2", Parsed{Weekday: "mon", Month: "jan", Day: "2", DateSep: " "}},
		{"date", "2 Jan", Parsed{Month: "jan", Day: "2", DateSep: " "}},
		{"date", "2 Jan 2006", Parsed{Month: "jan", Day: "2", Year: "2006", DateSep: " "}},
		{"date", "2 Jan 06", Parsed{Month: "jan", Day: "2", Year: "06", DateSep: " "}},

		{"time", "15:04:05", Parsed{Hour: "15", Minute: "04", Second: "05", TimeSep: ":"}},
		{"time", "15:04:05 pdt", Parsed{Hour: "15", Minute: "04", Second: "05", Zone: "pdt", Offset: "-0700", TimeSep: ":"}},
		{"time", "15:04:05.9999", Parsed{Hour: "15", Minute: "04", Second: "05", FracSecond: "9999", TimeSep: ":"}},
		{"time", "15:04", Parsed{Hour: "15", Minute: "04", TimeSep: ":"}},
		{"time", "3:04P.M.", Parsed{Hour: "3", Minute: "04", Period: "p", TimeSep: ":"}},
		{"time", "3:04am", Parsed{Hour: "3", Minute: "04", Period: "a", TimeSep: ":"}},
		{"time", "3:04am EST", Parsed{Hour: "3", Minute: "04", Period: "a", Zone: "est", Offset: "-0500", TimeSep: ":"}},
		{"time", "3:04am -0500", Parsed{Hour: "3", Minute: "04", Period: "a", Offset: "-0500", TimeSep: ":"}},
		{"time", "3:04am +05:00", Parsed{Hour: "3", Minute: "04", Period: "a", Offset: "+0500", TimeSep: ":"}},
		{"time", "3:04am -0500 EST", Parsed{Hour: "3", Minute: "04", Period: "a", Zone: "est", Offset: "-0500", TimeSep: ":"}},

		{"parse", "Mon Jan 1 2006 15:04:05 MST", Parsed{Weekday: "mon", Month: "jan", Day: "1", Year: "2006", Hour: "15", Minute: "04", Second: "05", Zone: "mst", Offset: "-0700", DateSep: " ", TimeSep: ":"}},
	}

	p := NewParser(locale.EnUS)
	p.Trace = true
	for _, test := range tests {
		// for _, test := range tests[len(tests)-1:] {
		t.Run(test.fn+":"+test.text, func(t *testing.T) {
			testValid(t, p, test.fn, test.text, test.parsed)
		})
		t.Run("parse:"+test.text, func(t *testing.T) {
			testValid(t, p, "parse", test.text, test.parsed)
		})
	}
}

func TestParserFrFR(t *testing.T) {
	tests := []struct {
		fn     string
		text   string
		parsed Parsed
	}{
		{"date", "2006-01-02", Parsed{Year: "2006", Month: "01", Day: "02", DateSep: "-"}},
		{"date", "2/1/2006", Parsed{Month: "1", Day: "2", Year: "2006", DateSep: "/"}},
		{"date", "2 janv 2006", Parsed{Month: "jan", Day: "2", Year: "2006", DateSep: " "}},
		{"date", "lundi, 2 janvier", Parsed{Weekday: "mon", Month: "jan", Day: "2", DateSep: " "}},

		{"time", "15:04:05,9999", Parsed{Hour: "15", Minute: "04", Second: "05", FracSecond: "9999", TimeSep: ":"}},
		{"time", "15 h 04", Parsed{Hour: "15", Minute: "04", HourSep: "h"}},
		{"time", "15h04", Parsed{Hour: "15", Minute: "04", HourSep: "h"}},

		{"parse", "lundi, 2/1/06 15:04:05,9999", Parsed{Weekday: "mon", Month: "1", Day: "2", Year: "06", Hour: "15", Minute: "04", Second: "05", FracSecond: "9999", DateSep: "/", TimeSep: ":"}},
	}

	p := NewParser(locale.FrFR)
	p.Trace = true
	for _, test := range tests {
		t.Run(test.fn+":"+test.text, func(t *testing.T) {
			testValid(t, p, test.fn, test.text, test.parsed)
		})
	}
}

func TestParserErrorEnUS(t *testing.T) {
	tests := []struct {
		fn   string
		text string
		err  string
	}{
		{"date", "2006", "invalid month"},
		{"time", "3:04am +1000 EST", "does not match given offset"},
	}

	p := NewParser(locale.EnUS)
	p.Trace = true
	for _, test := range tests {
		t.Run(test.fn+":"+test.text, func(t *testing.T) {
			check := func(have Parsed, err error) {
				if err == nil {
					t.Fatalf("expected error %v, have: %v", err, have)
				}
				if !strings.Contains(err.Error(), test.err) {
					t.Fatalf("\n have: %v want: %v\n", err.Error(), test.err)
				}
			}
			if test.fn == "date" {
				parsed, err := p.ParseDate(test.text)
				check(parsed, err)
			}
			if test.fn == "time" {
				parsed, err := p.ParseTime(test.text)
				check(parsed, err)
			}
		})
	}
}

func testValid(t *testing.T, p *Parser, fn string, text string, want Parsed) {
	check := func(have Parsed, want Parsed, err error) {
		if err != nil {
			t.Fatalf("unexpected error: %v \n have: %v \n tokens: %v", err, have, p.tokens)
		}
		if have != want {
			t.Errorf("\n have: %v \n want: %v", have, want)
		}
	}
	if fn == "date" {
		have, err := p.ParseDate(text)
		check(have, want, err)
	}
	if fn == "time" {
		have, err := p.ParseTime(text)
		check(have, want, err)
	}
	if fn == "parse" {
		have, err := p.Parse(text)
		check(have, want, err)
	}
}
