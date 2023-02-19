package main

import (
	"fmt"
	"log"
	"time"

	"github.com/blackchip-org/ptime"
	"github.com/blackchip-org/ptime/locale"
)

func main() {
	p, err := ptime.Parse(locale.EnUS, "3:04:05pm MST")
	if err != nil {
		log.Panic(err)
	}
	t, err := ptime.Time(p, time.Now())
	if err != nil {
		log.Panic(err)
	}
	f := ptime.Format(locale.EnUS, "[hour]:[minute]:[second] [offset]", t)
	fmt.Println(f)
}
