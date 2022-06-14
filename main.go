package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"text/template"
	"time"

	holiday "github.com/holiday-jp/holiday_jp-go"
	"github.com/pkg/errors"
)

type calendar struct {
	weekTemplate *template.Template
	target       time.Time
	month        *month
	buf          *strings.Builder
}

func newCalendar(target time.Time) (*calendar, error) {
	tmpl, err := template.New("week").Parse(weekTemplate)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &calendar{
		weekTemplate: tmpl,
		target:       target,
		month:        &month{},
		buf:          &strings.Builder{},
	}, nil
}

func (c *calendar) String() string {
	return c.buf.String()
}

const (
	weekTemplate = `|{{ if eq .Sunday.HolidayType 1 }} <font color="red">{{ if .Sunday.IsTargetMonth }}<b>{{ end }}{{.Sunday.N}}</font> {{ else if eq .Sunday.HolidayType 2 }} <font color="blue">{{ if .Sunday.IsTargetMonth }}<b>{{ end }}{{.Sunday.N}}</font> {{ else }} {{ if .Sunday.IsTargetMonth }}<b>{{ end }}{{.Sunday.N}} {{ end }}` +
		`|{{ if eq .Monday.HolidayType 1 }} <font color="red">{{ if .Monday.IsTargetMonth }}<b>{{ end }}{{.Monday.N}}</font> {{ else if eq .Monday.HolidayType 2 }} <font color="blue">{{ if .Monday.IsTargetMonth }}<b>{{ end }}{{.Monday.N}}</font> {{ else }} {{ if .Monday.IsTargetMonth }}<b>{{ end }}{{.Monday.N}} {{ end }}` +
		`|{{ if eq .Tuesday.HolidayType 1 }} <font color="red">{{ if .Tuesday.IsTargetMonth }}<b>{{ end }}{{.Tuesday.N}}</font> {{ else if eq .Tuesday.HolidayType 2 }} <font color="blue">{{ if .Tuesday.IsTargetMonth }}<b>{{ end }}{{.Tuesday.N}}</font> {{ else }} {{ if .Tuesday.IsTargetMonth }}<b>{{ end }}{{.Tuesday.N}} {{ end }}` +
		`|{{ if eq .Wednesday.HolidayType 1 }} <font color="red">{{ if .Wednesday.IsTargetMonth }}<b>{{ end }}{{.Wednesday.N}}</font> {{ else if eq .Wednesday.HolidayType 2 }} <font color="blue">{{ if .Wednesday.IsTargetMonth }}<b>{{ end }}{{.Wednesday.N}}</font> {{ else }} {{ if .Wednesday.IsTargetMonth }}<b>{{ end }}{{.Wednesday.N}} {{ end }}` +
		`|{{ if eq .Thursday.HolidayType 1 }} <font color="red">{{ if .Thursday.IsTargetMonth }}<b>{{ end }}{{.Thursday.N}}</font> {{ else if eq .Thursday.HolidayType 2 }} <font color="blue">{{ if .Thursday.IsTargetMonth }}<b>{{ end }}{{.Thursday.N}}</font> {{ else }} {{ if .Thursday.IsTargetMonth }}<b>{{ end }}{{.Thursday.N}} {{ end }}` +
		`|{{ if eq .Friday.HolidayType 1 }} <font color="red">{{ if .Friday.IsTargetMonth }}<b>{{ end }}{{.Friday.N}}</font> {{ else if eq .Friday.HolidayType 2 }} <font color="blue">{{ if .Friday.IsTargetMonth }}<b>{{ end }}{{.Friday.N}}</font> {{ else }} {{ if .Friday.IsTargetMonth }}<b>{{ end }}{{.Friday.N}} {{ end }}` +
		`|{{ if eq .Saturday.HolidayType 1 }} <font color="red">{{ if .Saturday.IsTargetMonth }}<b>{{ end }}{{.Saturday.N}}</font> {{ else if eq .Saturday.HolidayType 2 }} <font color="blue">{{ if .Saturday.IsTargetMonth }}<b>{{ end }}{{.Saturday.N}}</font> {{ else }} {{ if .Saturday.IsTargetMonth }}<b>{{ end }}{{.Saturday.N}} {{ end }}|`
)

type month struct {
	weeks []*week
}

type week struct {
	Sunday    *day
	Monday    *day
	Tuesday   *day
	Wednesday *day
	Thursday  *day
	Friday    *day
	Saturday  *day
}

type day struct {
	N             uint
	HolidayType   holidayType
	IsTargetMonth bool
}

type holidayType uint

const (
	notHoliday  holidayType = iota // 平日
	redHoliday                     // 日曜・祝日
	blueHoliday                    // 土曜
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
	flag.Parse()

	if retreat != 0 && advance != 0 {
		fmt.Fprintln(os.Stderr, errors.New("Please use either"))
		os.Exit(1)
	}

	targetDate := time.Now().In(jst)
	if retreat != 0 {
		targetDate = targetDate.AddDate(0, -retreat, 0)
	}
	if advance != 0 {
		targetDate = targetDate.AddDate(0, advance, 0)
	}

	calendar, err := newCalendar(targetDate)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if err := calendar.build(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Print(calendar)
}

func (c *calendar) build() error {
	if err := c.buildHeader(); err != nil {
		return errors.WithStack(err)
	}
	c.calculate()
	return errors.WithStack(c.render())
}

const (
	title     = `#### %d年%d月`
	header    = `|<font color="red">日</font>|月|火|水|木|金|<font color="blue">土</font>|`
	partition = `|--|--|--|--|--|--|--|`
)

func (c *calendar) buildHeader() error {
	_, err := c.buf.WriteString(fmt.Sprintf(title+"\n", c.target.Year(), c.target.Month()))
	if err != nil {
		return errors.WithStack(err)
	}
	_, err = c.buf.WriteString(header + "\n")
	if err != nil {
		return errors.WithStack(err)
	}
	_, err = c.buf.WriteString(partition + "\n")
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (c *calendar) calculate() {
	c.month.calculateWeeks(c.target)
}

func (c *calendar) render() error {
	for _, w := range c.month.weeks {
		if err := c.weekTemplate.Execute(c.buf, w); err != nil {
			return errors.WithStack(err)
		}
		_, err := c.buf.WriteString("\n")
		if err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

func (m *month) calculateWeeks(targetDate time.Time) {
	targetMonth := targetDate.Month()
	weeks := []*week{m.calculateWeek(targetDate, targetMonth)}

	retreat := targetDate
	for {
		retreat = retreat.AddDate(0, 0, -7)
		w := m.calculateWeek(retreat, targetMonth)
		weeks = append([]*week{w}, weeks...)
		if retreat.Month() != targetDate.Month() {
			break
		}
	}

	advance := targetDate
	for {
		advance = advance.AddDate(0, 0, 7)
		w := m.calculateWeek(advance, targetMonth)
		weeks = append(weeks, w)
		if advance.Month() != targetDate.Month() {
			break
		}
	}

	m.weeks = weeks
}

func (m *month) calculateWeek(date time.Time, targetMonth time.Month) *week {
	w := &week{}

	// dateの週の日曜日まで遡る
	retreat := date
	for {
		w.calculateDay(retreat, targetMonth)
		if retreat.Weekday() == time.Sunday {
			break
		}
		retreat = retreat.AddDate(0, 0, -1)
	}

	// dateの週の土曜日まで進む
	advance := date
	for {
		w.calculateDay(advance, targetMonth)
		if advance.Weekday() == time.Saturday {
			break
		}
		advance = advance.AddDate(0, 0, 1)
	}

	return w
}

func (w *week) calculateDay(date time.Time, targetMonth time.Month) {
	day := &day{
		N:             uint(date.Day()),
		IsTargetMonth: date.Month() == targetMonth,
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

func (d *day) calculateHoliday(date time.Time) {
	d.HolidayType = notHoliday
	switch {
	case holiday.IsHoliday(date):
		d.HolidayType = redHoliday
	case date.Weekday() == time.Sunday:
		d.HolidayType = redHoliday
	case date.Weekday() == time.Saturday:
		d.HolidayType = blueHoliday
	}
}
