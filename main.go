package main

import (
	"fmt"
	"os"
	"strings"
	"text/template"
	"time"

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
	now := time.Now()

	calendar, err := buildCalendar(now)
	if err != nil {
		fmt.Errorf("%+v", err)
		os.Exit(1)
	}

	fmt.Print(calendar)
}

const (
	calendarHeader = `|<font color="red">日</font>|月|火|水|木|金|<font color="blue">土</font>|`
	calendarSplit  = `|--|--|--|--|--|--|--|`
	weekTemplate   = `|<font color="red">{{.Sunday}}</font>|{{.Monday}}|{{.Tuesday}}|{{.Wednesday}}|{{.Thursday}}|{{.Friday}}|<font color="blue">{{.Saturday}}</font>|`
)

func buildCalendar(date time.Time) (string, error) {
	buf := &strings.Builder{}
	_, err := buf.WriteString(calendarHeader + "\n")
	if err != nil {
		return "", errors.WithStack(err)
	}
	_, err = buf.WriteString(calendarSplit + "\n")
	if err != nil {
		return "", errors.WithStack(err)
	}

	weeks := calculateWeeks(date)
	for _, wk := range weeks {
		tmpl, err := template.New("week").Parse(weekTemplate)
		if err != nil {
			return "", errors.WithStack(err)
		}
		if err := tmpl.Execute(buf, wk); err != nil {
			return "", errors.WithStack(err)
		}
		_, err = buf.WriteString("\n")
		if err != nil {
			return "", errors.WithStack(err)
		}
	}

	return buf.String(), nil
}

func calculateWeeks(date time.Time) []week {
	current := date
	wk := calculateWeek(current)
	weeks := []week{wk}

	retreat := current
	// currentの前月の最終週まで遡る
	for {
		retreat = retreat.AddDate(0, 0, -7)
		if retreat.Month() != current.Month() {
			wk := calculateWeek(retreat)
			weeks = append([]week{wk}, weeks...)
			break
		}
		wk := calculateWeek(retreat)
		weeks = append([]week{wk}, weeks...)
	}

	// currentの次月の初週まで進む
	advance := current
	for {
		advance = advance.AddDate(0, 0, 7)
		if advance.Month() != current.Month() {
			wk := calculateWeek(advance)
			weeks = append(weeks, wk)
			break
		}
		wk := calculateWeek(advance)
		weeks = append(weeks, wk)
	}

	return weeks
}

func calculateWeek(point time.Time) week {
	wk := week{}
	// pointの週の日曜日まで遡る
	retreat := point
	for {
		if retreat.Weekday() == time.Sunday {
			wk.Sunday = uint(retreat.Day())
			break
		}
		switch retreat.Weekday() {
		case time.Monday:
			wk.Monday = uint(retreat.Day())
		case time.Tuesday:
			wk.Tuesday = uint(retreat.Day())
		case time.Wednesday:
			wk.Wednesday = uint(retreat.Day())
		case time.Thursday:
			wk.Thursday = uint(retreat.Day())
		case time.Friday:
			wk.Friday = uint(retreat.Day())
		case time.Saturday:
			wk.Saturday = uint(retreat.Day())
		}
		retreat = retreat.AddDate(0, 0, -1)
	}

	// pointの週の土曜日まで進む
	advance := point
	for {
		if advance.Weekday() == time.Saturday {
			wk.Saturday = uint(advance.Day())
			break
		}
		switch advance.Weekday() {
		case time.Monday:
			wk.Monday = uint(advance.Day())
		case time.Tuesday:
			wk.Tuesday = uint(advance.Day())
		case time.Wednesday:
			wk.Wednesday = uint(advance.Day())
		case time.Thursday:
			wk.Thursday = uint(advance.Day())
		case time.Friday:
			wk.Friday = uint(advance.Day())
		}
		advance = advance.AddDate(0, 0, 1)
	}
	return wk
}
