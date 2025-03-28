package calma

import (
	"bytes"
	"fmt"
	"sync"
	"text/template"
	"time"

	"github.com/gomarkdown/markdown"
	holiday "github.com/holiday-jp/holiday_jp-go"
	"golang.org/x/sync/errgroup"
	"golang.org/x/xerrors"
)

type Calendar struct {
	weekTemplate *template.Template
	target       time.Time
	month        *month
	buf          *bytes.Buffer

	concurrent bool
}

func NewCalendar(target time.Time) (*Calendar, error) {
	tmpl, err := template.New("week").Parse(weekTemplate)
	if err != nil {
		return nil, xerrors.Errorf("failed to template.New: %w", err)
	}

	c := &Calendar{
		weekTemplate: tmpl,
		target:       target,
		month:        &month{},
		buf:          &bytes.Buffer{},
		concurrent:   false,
	}

	if err := c.build(); err != nil {
		return nil, xerrors.Errorf("failed to build calendar: %w", err)
	}

	return c, nil
}

func NewCalendarConcurrency(target time.Time) (*Calendar, error) {
	tmpl, err := template.New("week").Parse(weekTemplate)
	if err != nil {
		return nil, xerrors.Errorf("failed to template.New: %w", err)
	}

	c := &Calendar{
		weekTemplate: tmpl,
		target:       target,
		month:        &month{},
		buf:          &bytes.Buffer{},
		concurrent:   true,
	}

	if err := c.build(); err != nil {
		return nil, xerrors.Errorf("failed to build calendar: %w", err)
	}

	return c, nil
}

func (c *Calendar) String() string {
	return c.buf.String()
}

func (c *Calendar) HTML() string {
	md := c.buf.Bytes()
	html := markdown.ToHTML(md, nil, nil)
	return string(html)
}

const (
	weekTemplate = `{{ if eq .Sunday.HolidayType 1 }} <font color="red">{{ if .Sunday.IsTargetMonth }}<b>{{ end }}{{.Sunday.N}}</font> {{ else if eq .Sunday.HolidayType 2 }} <font color="blue">{{ if .Sunday.IsTargetMonth }}<b>{{ end }}{{.Sunday.N}}</font> {{ else }} {{ if .Sunday.IsTargetMonth }}<b>{{ end }}{{.Sunday.N}} {{ end }}` +
		`|{{ if eq .Monday.HolidayType 1 }} <font color="red">{{ if .Monday.IsTargetMonth }}<b>{{ end }}{{.Monday.N}}</font> {{ else if eq .Monday.HolidayType 2 }} <font color="blue">{{ if .Monday.IsTargetMonth }}<b>{{ end }}{{.Monday.N}}</font> {{ else }} {{ if .Monday.IsTargetMonth }}<b>{{ end }}{{.Monday.N}} {{ end }}` +
		`|{{ if eq .Tuesday.HolidayType 1 }} <font color="red">{{ if .Tuesday.IsTargetMonth }}<b>{{ end }}{{.Tuesday.N}}</font> {{ else if eq .Tuesday.HolidayType 2 }} <font color="blue">{{ if .Tuesday.IsTargetMonth }}<b>{{ end }}{{.Tuesday.N}}</font> {{ else }} {{ if .Tuesday.IsTargetMonth }}<b>{{ end }}{{.Tuesday.N}} {{ end }}` +
		`|{{ if eq .Wednesday.HolidayType 1 }} <font color="red">{{ if .Wednesday.IsTargetMonth }}<b>{{ end }}{{.Wednesday.N}}</font> {{ else if eq .Wednesday.HolidayType 2 }} <font color="blue">{{ if .Wednesday.IsTargetMonth }}<b>{{ end }}{{.Wednesday.N}}</font> {{ else }} {{ if .Wednesday.IsTargetMonth }}<b>{{ end }}{{.Wednesday.N}} {{ end }}` +
		`|{{ if eq .Thursday.HolidayType 1 }} <font color="red">{{ if .Thursday.IsTargetMonth }}<b>{{ end }}{{.Thursday.N}}</font> {{ else if eq .Thursday.HolidayType 2 }} <font color="blue">{{ if .Thursday.IsTargetMonth }}<b>{{ end }}{{.Thursday.N}}</font> {{ else }} {{ if .Thursday.IsTargetMonth }}<b>{{ end }}{{.Thursday.N}} {{ end }}` +
		`|{{ if eq .Friday.HolidayType 1 }} <font color="red">{{ if .Friday.IsTargetMonth }}<b>{{ end }}{{.Friday.N}}</font> {{ else if eq .Friday.HolidayType 2 }} <font color="blue">{{ if .Friday.IsTargetMonth }}<b>{{ end }}{{.Friday.N}}</font> {{ else }} {{ if .Friday.IsTargetMonth }}<b>{{ end }}{{.Friday.N}} {{ end }}` +
		`|{{ if eq .Saturday.HolidayType 1 }} <font color="red">{{ if .Saturday.IsTargetMonth }}<b>{{ end }}{{.Saturday.N}}</font> {{ else if eq .Saturday.HolidayType 2 }} <font color="blue">{{ if .Saturday.IsTargetMonth }}<b>{{ end }}{{.Saturday.N}}</font> {{ else }} {{ if .Saturday.IsTargetMonth }}<b>{{ end }}{{.Saturday.N}} {{ end }}`
)

type month struct {
	weeks []*week

	concurrent bool
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

func (c *Calendar) build() error {
	if err := c.buildHeader(); err != nil {
		return xerrors.Errorf("failed to buildHeader: %w", err)
	}

	c.calculate()

	if err := c.render(); err != nil {
		return xerrors.Errorf("failed to render: %w", err)
	}

	return nil
}

const (
	title     = `#### %d年%d月`
	header    = `<font color="red">日</font>|月|火|水|木|金|<font color="blue">土</font>`
	partition = `--------|--------|--------|--------|--------|--------|--------`
)

func (c *Calendar) buildHeader() error {
	const failedMsg = "failed to WriteString: %w"

	_, err := c.buf.WriteString(fmt.Sprintf(title+"\n", c.target.Year(), c.target.Month()))
	if err != nil {
		return xerrors.Errorf(failedMsg, err)
	}
	_, err = c.buf.WriteString(header + "\n")
	if err != nil {
		return xerrors.Errorf(failedMsg, err)
	}
	_, err = c.buf.WriteString(partition + "\n")
	if err != nil {
		return xerrors.Errorf(failedMsg, err)
	}

	return nil
}

func (c *Calendar) calculate() {
	c.month.concurrent = c.concurrent
	c.month.calculateWeeks(c.target)
}

func (c *Calendar) render() error {
	if c.concurrent {
		m := sync.Map{}
		eg := errgroup.Group{}
		for i, w := range c.month.weeks {
			i, w := i, w
			eg.Go(func() error {
				buf := &bytes.Buffer{}
				if err := c.weekTemplate.Execute(buf, w); err != nil {
					return xerrors.Errorf("failed to template.Execute: %w", err)
				}
				if _, err := buf.WriteString("\n"); err != nil {
					return xerrors.Errorf("failed to WriteString: %w", err)
				}

				m.Store(i, buf)
				return nil
			})
		}
		if err := eg.Wait(); err != nil {
			return xerrors.Errorf("failed: %w", err)
		}

		for i := range c.month.weeks {
			if v, ok := m.Load(i); ok {
				b := v.(*bytes.Buffer)
				if _, err := b.WriteTo(c.buf); err != nil {
					return xerrors.Errorf("failed to WriteTo: %w", err)
				}
			}
		}

		return nil
	}

	for _, w := range c.month.weeks {
		buf := &bytes.Buffer{}
		if err := c.weekTemplate.Execute(buf, w); err != nil {
			return xerrors.Errorf("failed to template.Execute: %w", err)
		}
		if _, err := buf.WriteString("\n"); err != nil {
			return xerrors.Errorf("failed to WriteString: %w", err)
		}
		if _, err := buf.WriteTo(c.buf); err != nil {
			return xerrors.Errorf("failed to WriteTo: %w", err)
		}
	}

	return nil
}

func (m *month) calculateWeeks(pointDate time.Time) {
	retreat := func(pointDate time.Time, retWeeks chan<- []*week) {
		retreat := pointDate
		pointMonth := pointDate.Month()
		weeks := []*week{m.calculateWeek(pointDate, pointMonth)}
		for {
			retreat = retreat.AddDate(0, 0, -7)
			w := m.calculateWeek(retreat, pointMonth)
			weeks = append([]*week{w}, weeks...)
			if retreat.Month() != pointMonth {
				break
			}
		}
		retWeeks <- weeks
	}

	advance := func(pointDate time.Time, advWeeks chan<- []*week) {
		advance := pointDate
		pointMonth := pointDate.Month()
		weeks := make([]*week, 0, 8)
		for {
			advance = advance.AddDate(0, 0, 7)
			w := m.calculateWeek(advance, pointMonth)
			weeks = append(weeks, w)
			if advance.Month() != pointMonth {
				break
			}
		}
		advWeeks <- weeks
	}

	if m.concurrent {
		ret := make(chan []*week, 1)
		go retreat(pointDate, ret)

		adv := make(chan []*week, 1)
		go advance(pointDate, adv)

		m.weeks = append(<-ret, <-adv...)
	} else {
		ret := make(chan []*week, 1)
		retreat(pointDate, ret)

		adv := make(chan []*week, 1)
		advance(pointDate, adv)

		m.weeks = append(<-ret, <-adv...)
	}
}

func (m *month) calculateWeek(point time.Time, pointMonth time.Month) *week {
	retreatToSunday := func(w *week, done chan<- struct{}) {
		retreat := point
		for {
			w.calculateDay(retreat, pointMonth)
			if retreat.Weekday() == time.Sunday {
				break
			}
			retreat = retreat.AddDate(0, 0, -1)
		}
		done <- struct{}{}
	}

	advanceToSaturday := func(w *week, done chan<- struct{}) {
		advance := point
		for {
			if advance.Weekday() == time.Saturday {
				break
			}
			advance = advance.AddDate(0, 0, 1)
			w.calculateDay(advance, pointMonth)
		}
		done <- struct{}{}
	}

	if m.concurrent {
		w, done := &week{}, make(chan struct{}, 2)
		go retreatToSunday(w, done)
		go advanceToSaturday(w, done)

		<-done
		<-done

		return w
	}

	w, done := &week{}, make(chan struct{}, 2)
	retreatToSunday(w, done)
	advanceToSaturday(w, done)

	<-done
	<-done

	return w
}

func (w *week) calculateDay(date time.Time, pointMonth time.Month) {
	day := &day{
		N:             uint(date.Day()),
		IsTargetMonth: date.Month() == pointMonth,
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
