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

func main() {
	flag.Usage = func() {
		usage := "This CLI outputs Japanese calendar in Markdown. It supports national holidays.\n\nUsage of %s:\n"
		fmt.Fprintf(flag.CommandLine.Output(), usage, os.Args[0])
		flag.PrintDefaults()
	}
	var retreat, advance int
	var concurrent bool
	flag.IntVar(&retreat, "r", 0, "Number of months to retreat")
	flag.IntVar(&advance, "a", 0, "Number of months to advance")
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

	if retreat != 0 && advance != 0 {
		fmt.Fprintln(os.Stderr, xerrors.New("Please use either"))
		os.Exit(1)
	}

	targetDate := time.Now().In(jst)
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
		html := calendar.HTML()
		fmt.Print(html)
		return
	}
	fmt.Print(calendar)
}
