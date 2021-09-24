package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/pkg/errors"
	// "github.com/yut-kt/goholiday"
	holiday "github.com/holiday-jp/holiday_jp-go"
)

type holidayType uint

const (
	NotHoliday  holidayType = iota // 平日
	RedHoliday                     // 日曜・祝日
	BlueHoliday                    // 土曜
)

type Day struct {
	N             uint
	HolidayType   holidayType
	IsTargetMonth bool
}

type week struct {
	targetMonth time.Month

	Sunday    *Day
	Monday    *Day
	Tuesday   *Day
	Wednesday *Day
	Thursday  *Day
	Friday    *Day
	Saturday  *Day
}

type month struct {
	targetMonth time.Month
	weeks       []*week
}

var jst = time.FixedZone("JST", +9*60*60)

func main() {
	var retreat, advance int
	flag.IntVar(&retreat, "r", 0, "Number of months to retreat")
	flag.IntVar(&advance, "a", 0, "Number of months to advance")
	flag.Parse()

	if retreat != 0 && advance != 0 {
		fmt.Errorf("%+v", errors.New("Please use either"))
		os.Exit(1)
	}

	date := time.Now().In(jst)
	if retreat != 0 {
		date = date.AddDate(0, -retreat, 0)
	}
	if advance != 0 {
		date = date.AddDate(0, advance, 0)
	}

	calendar, err := buildCalendar(date)
	if err != nil {
		fmt.Errorf("%+v", err)
		os.Exit(1)
	}

	fmt.Print(calendar)
}

const (
	calendarTitleJP  = `#### %d年%d月`
	calendarHeaderJP = `|<font color="red">日</font>|月|火|水|木|金|<font color="blue">土</font>|`
	calendarSplit    = `|--|--|--|--|--|--|--|`
	weekTemplate     = `|{{ if eq .Sunday.HolidayType 1 }} <font color="red">{{ if .Sunday.IsTargetMonth }}<b>{{ end }}{{.Sunday.N}}</font> {{ else if eq .Sunday.HolidayType 2 }} <font color="blue">{{ if .Sunday.IsTargetMonth }}<b>{{ end }}{{.Sunday.N}}</font> {{ else }} {{ if .Sunday.IsTargetMonth }}<b>{{ end }}{{.Sunday.N}} {{ end }}` +
		`|{{ if eq .Monday.HolidayType 1 }} <font color="red">{{ if .Monday.IsTargetMonth }}<b>{{ end }}{{.Monday.N}}</font> {{ else if eq .Monday.HolidayType 2 }} <font color="blue">{{ if .Monday.IsTargetMonth }}<b>{{ end }}{{.Monday.N}}</font> {{ else }} {{ if .Monday.IsTargetMonth }}<b>{{ end }}{{.Monday.N}} {{ end }}` +
		`|{{ if eq .Tuesday.HolidayType 1 }} <font color="red">{{ if .Tuesday.IsTargetMonth }}<b>{{ end }}{{.Tuesday.N}}</font> {{ else if eq .Tuesday.HolidayType 2 }} <font color="blue">{{ if .Tuesday.IsTargetMonth }}<b>{{ end }}{{.Tuesday.N}}</font> {{ else }} {{ if .Tuesday.IsTargetMonth }}<b>{{ end }}{{.Tuesday.N}} {{ end }}` +
		`|{{ if eq .Wednesday.HolidayType 1 }} <font color="red">{{ if .Wednesday.IsTargetMonth }}<b>{{ end }}{{.Wednesday.N}}</font> {{ else if eq .Wednesday.HolidayType 2 }} <font color="blue">{{ if .Wednesday.IsTargetMonth }}<b>{{ end }}{{.Wednesday.N}}</font> {{ else }} {{ if .Wednesday.IsTargetMonth }}<b>{{ end }}{{.Wednesday.N}} {{ end }}` +
		`|{{ if eq .Thursday.HolidayType 1 }} <font color="red">{{ if .Thursday.IsTargetMonth }}<b>{{ end }}{{.Thursday.N}}</font> {{ else if eq .Thursday.HolidayType 2 }} <font color="blue">{{ if .Thursday.IsTargetMonth }}<b>{{ end }}{{.Thursday.N}}</font> {{ else }} {{ if .Thursday.IsTargetMonth }}<b>{{ end }}{{.Thursday.N}} {{ end }}` +
		`|{{ if eq .Friday.HolidayType 1 }} <font color="red">{{ if .Friday.IsTargetMonth }}<b>{{ end }}{{.Friday.N}}</font> {{ else if eq .Friday.HolidayType 2 }} <font color="blue">{{ if .Friday.IsTargetMonth }}<b>{{ end }}{{.Friday.N}}</font> {{ else }} {{ if .Friday.IsTargetMonth }}<b>{{ end }}{{.Friday.N}} {{ end }}` +
		`|{{ if eq .Saturday.HolidayType 1 }} <font color="red">{{ if .Saturday.IsTargetMonth }}<b>{{ end }}{{.Saturday.N}}</font> {{ else if eq .Saturday.HolidayType 2 }} <font color="blue">{{ if .Saturday.IsTargetMonth }}<b>{{ end }}{{.Saturday.N}}</font> {{ else }} {{ if .Saturday.IsTargetMonth }}<b>{{ end }}{{.Saturday.N}} {{ end }}|`
)

func buildCalendar(date time.Time) (string, error) {
	buf := &strings.Builder{}
	_, err := buf.WriteString(fmt.Sprintf(calendarTitleJP+"\n", date.Year(), date.Month()))
	if err != nil {
		return "", errors.WithStack(err)
	}
	_, err = buf.WriteString(calendarHeaderJP + "\n")
	if err != nil {
		return "", errors.WithStack(err)
	}
	_, err = buf.WriteString(calendarSplit + "\n")
	if err != nil {
		return "", errors.WithStack(err)
	}

	tmpl, err := template.New("week").Parse(weekTemplate)
	if err != nil {
		return "", errors.WithStack(err)
	}
	m := &month{targetMonth: date.Month()}
	m.calculateWeeks(date)
	for _, w := range m.weeks {
		if err := tmpl.Execute(buf, w); err != nil {
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
	w := m.calculateWeek(date)
	weeks := []*week{w}

	// dateの前月の最終週まで遡る
	retreat := date
	for {
		retreat = retreat.AddDate(0, 0, -7)
		w := m.calculateWeek(retreat)
		weeks = append([]*week{w}, weeks...)
		if retreat.Month() != date.Month() {
			break
		}
	}

	// dateの次月の初週まで進む
	advance := date
	for {
		advance = advance.AddDate(0, 0, 7)
		w := m.calculateWeek(advance)
		weeks = append(weeks, w)
		if advance.Month() != date.Month() {
			break
		}
	}

	m.weeks = weeks
}

func (m *month) calculateWeek(date time.Time) *week {
	w := &week{targetMonth: m.targetMonth}
	// dateの週の日曜日まで遡る
	retreat := date
	for {
		w.calculateDay(retreat)
		if retreat.Weekday() == time.Sunday {
			break
		}
		retreat = retreat.AddDate(0, 0, -1)
	}

	// dateの週の土曜日まで進む
	advance := date
	for {
		w.calculateDay(advance)
		if advance.Weekday() == time.Saturday {
			break
		}
		advance = advance.AddDate(0, 0, 1)
	}
	return w
}

func (w *week) calculateDay(date time.Time) {
	day := &Day{
		N:             uint(date.Day()),
		IsTargetMonth: date.Month() == w.targetMonth,
	}
	day.calculateHoliday(date)

	switch date.Weekday() {
	case time.Sunday:
		w.Sunday = day
	case time.Monday:
		w.Monday = day
	case time.Tuesday:
		w.Tuesday = day
	case time.Wednesday:
		w.Wednesday = day
	case time.Thursday:
		w.Thursday = day
	case time.Friday:
		w.Friday = day
	case time.Saturday:
		w.Saturday = day
	}
}

func (d *Day) calculateHoliday(date time.Time) {
	var hType holidayType
	switch {
	case holiday.IsHoliday(date):
		hType = RedHoliday
	case date.Weekday() == time.Sunday:
		hType = RedHoliday
	case date.Weekday() == time.Saturday:
		hType = BlueHoliday
	default:
		hType = NotHoliday
	}
	d.HolidayType = hType
}
