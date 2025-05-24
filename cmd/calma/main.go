package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"golang.org/x/xerrors"

	"github.com/ddddddO/calma"
)

var jst = time.FixedZone("JST", +9*60*60)

var (
	Version  = "unset"
	Revision = "unset"
)

type monthText struct {
	target time.Time
}

const (
	monthFormat = "2006-01"
	middleDay   = 14 // 月の中間あたりということで
)

func (m *monthText) UnmarshalText(text []byte) error {
	origin, err := time.Parse(monthFormat, string(text))
	if err != nil {
		return err
	}

	m.target = time.Date(
		origin.Year(),
		origin.Month(),
		middleDay,
		origin.Hour(),
		origin.Minute(),
		origin.Second(),
		origin.Nanosecond(),
		origin.Location(),
	)
	return nil
}

func (m *monthText) MarshalText() (text []byte, err error) {
	return []byte(m.target.Format(monthFormat)), nil
}

func main() {
	flag.Usage = func() {
		usage := "This CLI outputs Japanese calendar in Markdown. It supports national holidays.\n\nUsage of %s:\n"
		fmt.Fprintf(flag.CommandLine.Output(), usage, os.Args[0])
		flag.PrintDefaults()
	}

	now := time.Now().In(jst)
	target := monthText{target: time.Date(now.Year(), now.Month(), middleDay, now.Hour(), now.Minute(), now.Second(), now.Nanosecond(), now.Location())}
	flag.TextVar(&target, "t", &target, "Specify month ('2006-01')")

	var retreat, advance int
	flag.IntVar(&retreat, "r", 0, "Number of months to retreat")
	flag.IntVar(&advance, "a", 0, "Number of months to advance")

	var concurrent bool
	flag.BoolVar(&concurrent, "concurrent", false, "Concurrent processing (performance deteriorates)")
	var isHTML bool
	flag.BoolVar(&isHTML, "html", false, "Output html")
	var isVersion bool
	flag.BoolVar(&isVersion, "version", false, "Show the version")

	flag.Parse()

	if isVersion {
		fmt.Printf("calma version %s / revision %s\n", Version, Revision)
		os.Exit(0)
	}

	targetDate := target.target
	if retreat != 0 {
		targetDate = targetDate.AddDate(0, -retreat, 0)
	}
	if advance != 0 {
		targetDate = targetDate.AddDate(0, advance, 0)
	}

	var calendar *calma.Calendar
	var err error
	if concurrent {
		calendar, err = calma.NewCalendarConcurrency(targetDate)
	} else {
		calendar, err = calma.NewCalendar(targetDate)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, xerrors.Errorf("%+v", err))
		os.Exit(1)
	}

	if isHTML {
		fmt.Print(calendar.HTML())
	} else {
		fmt.Print(calendar)
	}
}
