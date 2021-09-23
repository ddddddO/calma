package main

import (
	"fmt"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/pkg/errors"
)

const (
	NotHoliday  = iota // 平日
	RedHoliday         // 日曜・祝日
	BlueHoliday        // 土曜
)

type Day struct {
	N           uint
	HolidayType uint
	IsThisMonth bool
}

type week struct {
	Sunday    *Day
	Monday    *Day
	Tuesday   *Day
	Wednesday *Day
	Thursday  *Day
	Friday    *Day
	Saturday  *Day
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
	weekTemplate   = `|<font color="red">{{.Sunday.N}}</font>|{{.Monday.N}}|{{.Tuesday.N}}|{{.Wednesday.N}}|{{.Thursday.N}}|{{.Friday.N}}|<font color="blue">{{.Saturday.N}}</font>|`
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

func calculateWeeks(date time.Time) []*week {
	current := date
	wk := calculateWeek(current)
	weeks := []*week{wk}

	retreat := current
	// currentの前月の最終週まで遡る
	for {
		retreat = retreat.AddDate(0, 0, -7)
		if retreat.Month() != current.Month() {
			wk := calculateWeek(retreat)
			weeks = append([]*week{wk}, weeks...)
			break
		}
		wk := calculateWeek(retreat)
		weeks = append([]*week{wk}, weeks...)
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

func calculateWeek(point time.Time) *week {
	wk := &week{}
	// pointの週の日曜日まで遡る
	retreat := point
	for {
		wk.calculateDay(retreat)
		if retreat.Weekday() == time.Sunday {
			break
		}
		retreat = retreat.AddDate(0, 0, -1)
	}

	// pointの週の土曜日まで進む
	advance := point
	for {
		wk.calculateDay(advance)
		if advance.Weekday() == time.Saturday {
			break
		}
		advance = advance.AddDate(0, 0, 1)
	}
	return wk
}

func (wk *week) calculateDay(date time.Time) {
	day := &Day{N: uint(date.Day())}
	switch date.Weekday() {
	case time.Sunday:
		day.HolidayType = RedHoliday
		wk.Sunday = day
	case time.Monday:
		wk.Monday = day
	case time.Tuesday:
		wk.Tuesday = day
	case time.Wednesday:
		wk.Wednesday = day
	case time.Thursday:
		wk.Thursday = day
	case time.Friday:
		wk.Friday = day
	case time.Saturday:
		day.HolidayType = BlueHoliday
		wk.Saturday = day
	}
}
