package main

import (
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/pkg/errors"
)

type week struct {
	Sunday    uint
	Monday    uint
	Tuesday   uint
	Wednesday uint
	Thursday  uint
	Friday    uint
	Saturday  uint
}

func main() {
	calendar, err := buildCalendar()
	if err != nil {
		fmt.Errorf("%+v", err)
		os.Exit(1)
	}

	fmt.Println(calendar)
}

const (
	calendarHeader = `|<font color="red">日</font>|月|火|水|木|金|<font color="blue">土</font>|`
	calendarSplit  = `|--|--|--|--|--|--|--|`
	weekTemplate   = `|{{.Sunday}}|{{.Monday}}|{{.Tuesday}}|{{.Wednesday}}|{{.Thursday}}|{{.Friday}}|{{.Saturday}}|`
)

func buildCalendar() (string, error) {
	buf := &strings.Builder{}
	_, err := buf.WriteString(calendarHeader + "\n")
	if err != nil {
		return "", errors.WithStack(err)
	}
	_, err = buf.WriteString(calendarSplit + "\n")
	if err != nil {
		return "", errors.WithStack(err)
	}

	wk := week{Sunday: 1, Monday: 2, Tuesday: 3, Wednesday: 4, Thursday: 5, Friday: 6, Saturday: 7}
	tmpl, err := template.New("week").Parse(weekTemplate)
	if err != nil {
		return "", errors.WithStack(err)
	}
	if err := tmpl.Execute(buf, wk); err != nil {
		return "", errors.WithStack(err)
	}

	return buf.String(), nil
}
