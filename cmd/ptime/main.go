package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/blackchip-org/ptime"
	"github.com/blackchip-org/ptime/locale"
)

var (
	dateOnly   bool
	format     string
	localeName string
	timeOnly   bool
	verbose    bool
)

func main() {
	log.SetFlags(0)
	flag.BoolVar(&dateOnly, "d", false, "only parse date")
	flag.StringVar(&format, "f", "", "format the result")
	flag.StringVar(&localeName, "l", "en-US", "set locale")
	flag.BoolVar(&timeOnly, "t", false, "only parse time")
	flag.BoolVar(&verbose, "v", false, "verbose")

	flag.Parse()

	text := strings.Join(flag.Args(), " ")
	l, ok := locale.Lookup(localeName)
	if !ok {
		log.Fatalf("locale '%v' not found", localeName)
	}

	p := ptime.NewParser(l)
	if verbose {
		p.Trace = true
	}

	var parseFn func(string) (ptime.Parsed, error)
	switch {
	case dateOnly:
		parseFn = p.ParseDate
	case timeOnly:
		parseFn = p.ParseTime
	default:
		parseFn = p.Parse
	}
	res, err := parseFn(text)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	if format != "" {
		t, err := ptime.Time(res, time.Now())
		if err != nil {
			log.Fatalf("unexpected error: %v", err)
		}
		fmt.Println(ptime.Format(l, format, t))
	} else {
		b, err := json.MarshalIndent(res, "", "  ")
		if err != nil {
			log.Panic(err)
		}
		fmt.Println(string(b))
	}
}
