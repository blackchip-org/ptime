# ptime

Date and time parser.

Experimental and a work in progress.

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

## License

MIT

## Feedback

Contact me at zc@blackchip.org



