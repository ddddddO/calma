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
	Sunday    Day
	Monday    Day
	Tuesday   Day
	Wednesday Day
	Thursday  Day
	Friday    Day
	Saturday  Day
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
			wk.Sunday.N = uint(retreat.Day())
			wk.Sunday.HolidayType = RedHoliday
			break
		}
		switch retreat.Weekday() {
		case time.Monday:
			wk.Monday.N = uint(retreat.Day())
		case time.Tuesday:
			wk.Tuesday.N = uint(retreat.Day())
		case time.Wednesday:
			wk.Wednesday.N = uint(retreat.Day())
		case time.Thursday:
			wk.Thursday.N = uint(retreat.Day())
		case time.Friday:
			wk.Friday.N = uint(retreat.Day())
		case time.Saturday:
			wk.Saturday.N = uint(retreat.Day())
			wk.Saturday.HolidayType = BlueHoliday
		}
		retreat = retreat.AddDate(0, 0, -1)
	}

	// pointの週の土曜日まで進む
	advance := point
	for {
		if advance.Weekday() == time.Saturday {
			wk.Saturday.N = uint(advance.Day())
			wk.Saturday.HolidayType = BlueHoliday
			break
		}
		switch advance.Weekday() {
		case time.Monday:
			wk.Monday.N = uint(advance.Day())
		case time.Tuesday:
			wk.Tuesday.N = uint(advance.Day())
		case time.Wednesday:
			wk.Wednesday.N = uint(advance.Day())
		case time.Thursday:
			wk.Thursday.N = uint(advance.Day())
		case time.Friday:
			wk.Friday.N = uint(advance.Day())
		}
		advance = advance.AddDate(0, 0, 1)
	}
	return wk
}
