package main

import (
	"fmt"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/pkg/errors"
	"github.com/yut-kt/goholiday"
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
	executedMonth time.Month

	Sunday    *Day
	Monday    *Day
	Tuesday   *Day
	Wednesday *Day
	Thursday  *Day
	Friday    *Day
	Saturday  *Day
}

type month struct {
	executedMonth time.Month
	weeks         []*week
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
	weekTemplate   = `|{{ if eq .Sunday.HolidayType 1 }} <font color="red">{{ if .Sunday.IsThisMonth }}<b>{{ end }}{{.Sunday.N}}</font> {{ else if eq .Sunday.HolidayType 2 }} <font color="blue">{{ if .Sunday.IsThisMonth }}<b>{{ end }}{{.Sunday.N}}</font> {{ else }} {{ if .Sunday.IsThisMonth }}<b>{{ end }}{{.Sunday.N}} {{ end }}` +
		`|{{ if eq .Monday.HolidayType 1 }} <font color="red">{{ if .Monday.IsThisMonth }}<b>{{ end }}{{.Monday.N}}</font> {{ else if eq .Monday.HolidayType 2 }} <font color="blue">{{ if .Monday.IsThisMonth }}<b>{{ end }}{{.Monday.N}}</font> {{ else }} {{ if .Monday.IsThisMonth }}<b>{{ end }}{{.Monday.N}} {{ end }}` +
		`|{{ if eq .Tuesday.HolidayType 1 }} <font color="red">{{ if .Tuesday.IsThisMonth }}<b>{{ end }}{{.Tuesday.N}}</font> {{ else if eq .Tuesday.HolidayType 2 }} <font color="blue">{{ if .Tuesday.IsThisMonth }}<b>{{ end }}{{.Tuesday.N}}</font> {{ else }} {{ if .Tuesday.IsThisMonth }}<b>{{ end }}{{.Tuesday.N}} {{ end }}` +
		`|{{ if eq .Wednesday.HolidayType 1 }} <font color="red">{{ if .Wednesday.IsThisMonth }}<b>{{ end }}{{.Wednesday.N}}</font> {{ else if eq .Wednesday.HolidayType 2 }} <font color="blue">{{ if .Wednesday.IsThisMonth }}<b>{{ end }}{{.Wednesday.N}}</font> {{ else }} {{ if .Wednesday.IsThisMonth }}<b>{{ end }}{{.Wednesday.N}} {{ end }}` +
		`|{{ if eq .Thursday.HolidayType 1 }} <font color="red">{{ if .Thursday.IsThisMonth }}<b>{{ end }}{{.Thursday.N}}</font> {{ else if eq .Thursday.HolidayType 2 }} <font color="blue">{{ if .Thursday.IsThisMonth }}<b>{{ end }}{{.Thursday.N}}</font> {{ else }} {{ if .Thursday.IsThisMonth }}<b>{{ end }}{{.Thursday.N}} {{ end }}` +
		`|{{ if eq .Friday.HolidayType 1 }} <font color="red">{{ if .Friday.IsThisMonth }}<b>{{ end }}{{.Friday.N}}</font> {{ else if eq .Friday.HolidayType 2 }} <font color="blue">{{ if .Friday.IsThisMonth }}<b>{{ end }}{{.Friday.N}}</font> {{ else }} {{ if .Friday.IsThisMonth }}<b>{{ end }}{{.Friday.N}} {{ end }}` +
		`|{{ if eq .Saturday.HolidayType 1 }} <font color="red">{{ if .Saturday.IsThisMonth }}<b>{{ end }}{{.Saturday.N}}</font> {{ else if eq .Saturday.HolidayType 2 }} <font color="blue">{{ if .Saturday.IsThisMonth }}<b>{{ end }}{{.Saturday.N}}</font> {{ else }} {{ if .Saturday.IsThisMonth }}<b>{{ end }}{{.Saturday.N}} {{ end }}|`
)

func buildCalendar(date time.Time) (string, error) {
	buf := &strings.Builder{}
	_, err := buf.WriteString(fmt.Sprintf("#### %d年%d月\n", date.Year(), date.Month()))
	if err != nil {
		return "", errors.WithStack(err)
	}
	_, err = buf.WriteString(calendarHeader + "\n")
	if err != nil {
		return "", errors.WithStack(err)
	}
	_, err = buf.WriteString(calendarSplit + "\n")
	if err != nil {
		return "", errors.WithStack(err)
	}

	m := &month{executedMonth: date.Month()}
	m.calculateWeeks(date)
	for _, wk := range m.weeks {
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

func (m *month) calculateWeeks(date time.Time) {
	current := date
	wk := m.calculateWeek(current)
	weeks := []*week{wk}

	retreat := current
	// currentの前月の最終週まで遡る
	for {
		retreat = retreat.AddDate(0, 0, -7)
		if retreat.Month() != current.Month() {
			wk := m.calculateWeek(retreat)
			weeks = append([]*week{wk}, weeks...)
			break
		}
		wk := m.calculateWeek(retreat)
		weeks = append([]*week{wk}, weeks...)
	}

	// currentの次月の初週まで進む
	advance := current
	for {
		advance = advance.AddDate(0, 0, 7)
		if advance.Month() != current.Month() {
			wk := m.calculateWeek(advance)
			weeks = append(weeks, wk)
			break
		}
		wk := m.calculateWeek(advance)
		weeks = append(weeks, wk)
	}

	m.weeks = weeks
}

func (m *month) calculateWeek(point time.Time) *week {
	wk := &week{executedMonth: m.executedMonth}
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
	day := &Day{
		N:           uint(date.Day()),
		IsThisMonth: date.Month() == wk.executedMonth,
		HolidayType: calculateHoliday(date),
	}

	switch date.Weekday() {
	case time.Sunday:
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
		wk.Saturday = day
	}
}

func calculateHoliday(date time.Time) uint {
	if goholiday.IsNationalHoliday(date) {
		return RedHoliday
	}
	if date.Weekday() == time.Sunday {
		return RedHoliday
	}
	if date.Weekday() == time.Saturday {
		return BlueHoliday
	}
	return NotHoliday
}
