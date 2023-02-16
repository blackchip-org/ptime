package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/blackchip-org/ptime"
	"github.com/blackchip-org/ptime/locale"
)

var (
	localeName string
	trace      bool
)

func main() {
	log.SetFlags(0)
	flag.StringVar(&localeName, "l", "en-US", "set locale")
	flag.BoolVar(&trace, "t", false, "trace")

	flag.Parse()

	text := strings.Join(flag.Args(), " ")
	l, ok := locale.Lookup(localeName)
	if !ok {
		log.Fatalf("locale '%v' not found", localeName)
	}

	p := ptime.NewParser(l)
	if trace {
		p.Trace = true
	}
	res, err := p.Parse(text)
	if err != nil {
		fmt.Println("error")
		log.Fatal(err)
	}
	b, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		log.Panic(err)
	}
	fmt.Println(string(b))
}
