# ptime

Parse date and times without knowing the layout ahead of time.

Experimental and a work in progress. More details to come later.

## Installation

Install [go](https://go.dev/dl/).

Install with:

    go install github.com/blackchip-org/ptime/cmd/ptime@latest

## Examples

    ptime -l en-US Mon Jan 2 2006 3:04:05pm MST

Output:

```json
{
  "Weekday": "Mon",
  "Year": "2006",
  "Month": "Jan",
  "Day": "2",
  "Hour": "3",
  "Minute": "04",
  "Second": "05",
  "Period": "PM",
  "Zone": "MST",
  "Offset": "-0700",
  "DateSep": " ",
  "TimeSep": ":"
}
```

    ptime -l fr-FR lundi, 2/1/06 15:04:05,9999

Output:

```json
{
  "Weekday": "Mon",
  "Year": "06",
  "Month": "1",
  "Day": "2",
  "Hour": "15",
  "Minute": "04",
  "Second": "05",
  "FracSecond": "9999",
  "DateSep": "/",
  "TimeSep": ":"
}
```

The only locales supported at the moment are
[en-US](https://github.com/blackchip-org/ptime/blob/main/locale/en.go) and
[fr-FR](https://github.com/blackchip-org/ptime/blob/main/locale/fr.go). The
[CLDR](https://cldr.unicode.org/) may be included at some point.

Supported layouts are shown in the tests here:

https://github.com/blackchip-org/ptime/blob/main/parser_test.go

The behavior of other layouts is undefined. Use the `-v` option on the
command line to get insight into the parsing process.

## License

MIT

## Feedback

Contact me at zc@blackchip.org



