package ptime

import (
	"fmt"
	"strconv"
	"time"

	"github.com/blackchip-org/ptime/locale"
)

func Time(p Parsed, now time.Time) (time.Time, error) {
	var year, mon, day, hour, min, sec, nsec int
	var loc *time.Location
	var err error

	if p.Year != "" {
		year, err = strconv.Atoi(p.Year)
		if err != nil {
			return time.Time{}, fmt.Errorf("invalid year: %v", p.Year)
		}
		if len(p.Year) == 2 {
			year = now.Year()/1000*1000 + year
		}
	} else {
		year = now.Year()
	}

	if p.Month != "" {
		m, ok := locale.EnUS.MonthNum[p.Month]
		if ok {
			mon = m
		}
		if !ok {
			mon, err = strconv.Atoi(p.Month)
			if err != nil {
				return time.Time{}, fmt.Errorf("invalid month: %v", p.Month)
			}
		}
	} else {
		mon = int(now.Month())
	}

	if p.Day != "" {
		day, err = strconv.Atoi(p.Day)
		if err != nil {
			return time.Time{}, fmt.Errorf("invalid day: %v", p.Month)
		}

		if len(p.Day) == 3 {
			if p.Month != "" {
				return time.Time{}, fmt.Errorf("must use either ordinal day or month")
			}
			mon = 1
		}
	} else {
		if p.Year == "" && p.Month == "" {
			day = now.Day()
		} else {
			day = 1
		}
	}

	if p.Hour != "" {
		hour, err = strconv.Atoi(p.Hour)
		if err != nil {
			return time.Time{}, fmt.Errorf("invalid hour: %v", p.Hour)
		}

		if p.Period != "" {
			num, ok := locale.EnUS.PeriodNum[p.Period]
			if !ok {
				return time.Time{}, fmt.Errorf("invalid period: %v", p.Period)
			}
			if num == int(locale.PM) || num == int(locale.Midnight) {
				hour += 12
			}
		}
	}

	if p.Minute != "" {
		min, err = strconv.Atoi(p.Minute)
		if err != nil {
			return time.Time{}, fmt.Errorf("invalid minute: %v", p.Minute)
		}
	}
	if p.Second != "" {
		sec, err = strconv.Atoi(p.Second)
		if err != nil {
			return time.Time{}, fmt.Errorf("invalid second: %v", p.Second)
		}
	}
	if p.FracSecond != "" {
		fsec, err := strconv.Atoi(p.FracSecond)
		if err != nil {
			return time.Time{}, fmt.Errorf("invalid fractional second: %v", p.FracSecond)
		}
		nsec = fsecToNsec(fsec)
	}

	if p.Offset != "" {
		o, err := strconv.Atoi(p.Offset)
		if err != nil {
			return time.Time{}, fmt.Errorf("invalid offset: %v", p.Offset)
		}
		oh := o / 100
		om := o % 100
		offset := (oh * 3600) + (om * 60)
		loc = time.FixedZone(p.Zone, offset)
	} else {
		loc = now.Location()
	}

	return time.Date(year, time.Month(mon), day, hour, min, sec, nsec, loc), nil
}

// There must be a better way to do this
func fsecToNsec(fsec int) int {
	sec, err := strconv.ParseFloat(fmt.Sprintf(".%v", fsec), 64)
	if err != nil {
		panic(err)
	}
	return int(sec * 1e9)
}

// func Format(locale locale.Locale, layout string, t time.Time) string {
// 	var token strings.Builder
// 	var result strings.Builder
// 	var inToken bool

// 	for _, ch := range layout {
// 		if !inToken {
// 			if ch == '[' {
// 				inToken = true
// 			} else {
// 				result.WriteRune(ch)
// 			}
// 			continue
// 		}
// 		if ch != ']' {
// 			token.WriteRune(ch)
// 			continue
// 		}
// 		field := token.String()
// 		token.Reset()
// 		inToken = false

// 		switch field {
// 		case "year":
// 			result.WriteString(fmt.Sprintf("%04d", t.Year()))
// 		case "year2":
// 			result.WriteString(fmt.Sprintf("%02d", t.Year() / 100))
// 		case "month":
// 			result.WriteString(fmt.Sprintf("%d", t.Month()))
// 		case "month2":
// 			result.WriteString(fmt.Sprintf("%02d", t.Month()))
// 		case "day":
// 			result.WriteString(fmt.Sprintf("%d", t.Day()))
// 		case "day2":
// 			result.WriteString(fmt.Sprintf("%02d", t.Day()))
// 		case "day3":
// 			result.WriteString(t.Format("002"))

// 	}
// }