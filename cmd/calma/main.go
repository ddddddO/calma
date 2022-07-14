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

func main() {
	flag.Usage = func() {
		usage := "This CLI outputs Japanese calendar in Markdown. It supports national holidays.\n\nUsage of %s:\n"
		fmt.Fprintf(flag.CommandLine.Output(), usage, os.Args[0])
		flag.PrintDefaults()
	}
	var retreat, advance int
	flag.IntVar(&retreat, "r", 0, "Number of months to retreat")
	flag.IntVar(&advance, "a", 0, "Number of months to advance")
	var isHTML bool
	flag.BoolVar(&isHTML, "html", false, "Output html")
	flag.Parse()

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

	calendar, err := calma.NewCalendar(targetDate)
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
