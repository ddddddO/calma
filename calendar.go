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

const (
	weekTemplate = `{{ if eq .Sunday.HolidayType 1 }} <font color="red">{{ if .Sunday.IsTargetMonth }}<b>{{ end }}{{.Sunday.N}}</font> {{ else if eq .Sunday.HolidayType 2 }} <font color="blue">{{ if .Sunday.IsTargetMonth }}<b>{{ end }}{{.Sunday.N}}</font> {{ else }} {{ if .Sunday.IsTargetMonth }}<b>{{ end }}{{.Sunday.N}} {{ end }}` +
		`|{{ if eq .Monday.HolidayType 1 }} <font color="red">{{ if .Monday.IsTargetMonth }}<b>{{ end }}{{.Monday.N}}</font> {{ else if eq .Monday.HolidayType 2 }} <font color="blue">{{ if .Monday.IsTargetMonth }}<b>{{ end }}{{.Monday.N}}</font> {{ else }} {{ if .Monday.IsTargetMonth }}<b>{{ end }}{{.Monday.N}} {{ end }}` +
		`|{{ if eq .Tuesday.HolidayType 1 }} <font color="red">{{ if .Tuesday.IsTargetMonth }}<b>{{ end }}{{.Tuesday.N}}</font> {{ else if eq .Tuesday.HolidayType 2 }} <font color="blue">{{ if .Tuesday.IsTargetMonth }}<b>{{ end }}{{.Tuesday.N}}</font> {{ else }} {{ if .Tuesday.IsTargetMonth }}<b>{{ end }}{{.Tuesday.N}} {{ end }}` +
		`|{{ if eq .Wednesday.HolidayType 1 }} <font color="red">{{ if .Wednesday.IsTargetMonth }}<b>{{ end }}{{.Wednesday.N}}</font> {{ else if eq .Wednesday.HolidayType 2 }} <font color="blue">{{ if .Wednesday.IsTargetMonth }}<b>{{ end }}{{.Wednesday.N}}</font> {{ else }} {{ if .Wednesday.IsTargetMonth }}<b>{{ end }}{{.Wednesday.N}} {{ end }}` +
		`|{{ if eq .Thursday.HolidayType 1 }} <font color="red">{{ if .Thursday.IsTargetMonth }}<b>{{ end }}{{.Thursday.N}}</font> {{ else if eq .Thursday.HolidayType 2 }} <font color="blue">{{ if .Thursday.IsTargetMonth }}<b>{{ end }}{{.Thursday.N}}</font> {{ else }} {{ if .Thursday.IsTargetMonth }}<b>{{ end }}{{.Thursday.N}} {{ end }}` +
		`|{{ if eq .Friday.HolidayType 1 }} <font color="red">{{ if .Friday.IsTargetMonth }}<b>{{ end }}{{.Friday.N}}</font> {{ else if eq .Friday.HolidayType 2 }} <font color="blue">{{ if .Friday.IsTargetMonth }}<b>{{ end }}{{.Friday.N}}</font> {{ else }} {{ if .Friday.IsTargetMonth }}<b>{{ end }}{{.Friday.N}} {{ end }}` +
		`|{{ if eq .Saturday.HolidayType 1 }} <font color="red">{{ if .Saturday.IsTargetMonth }}<b>{{ end }}{{.Saturday.N}}</font> {{ else if eq .Saturday.HolidayType 2 }} <font color="blue">{{ if .Saturday.IsTargetMonth }}<b>{{ end }}{{.Saturday.N}}</font> {{ else }} {{ if .Saturday.IsTargetMonth }}<b>{{ end }}{{.Saturday.N}} {{ end }}`
)

type calendar struct {
	weekTemplate *template.Template
	target       time.Time
	month        *month
	buf          *bytes.Buffer
}

func NewCalendar(target time.Time) (*calendar, error) {
	tmpl, err := template.New("week").Parse(weekTemplate)
	if err != nil {
		return nil, xerrors.Errorf("failed to template.New: %w", err)
	}

	c := &calendar{
		weekTemplate: tmpl,
		target:       target,
		month:        &month{},
		buf:          &bytes.Buffer{},
	}

	if err := c.build(); err != nil {
		return nil, xerrors.Errorf("failed to build calendar: %w", err)
	}

	return c, nil
}

func (c *calendar) String() string {
	return c.buf.String()
}

func (c *calendar) HTML() string {
	md := []byte(c.String())
	html := markdown.ToHTML(md, nil, nil)
	return string(html)
}

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

func (c *calendar) build() error {
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

func (c *calendar) buildHeader() error {
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

func (c *calendar) calculate() {
	c.month.calculateWeeks(c.target)
}

func (c *calendar) render() error {
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

func (m *month) calculateWeeks(targetDate time.Time) {
	targetMonth := targetDate.Month()

	retWeeks := make(chan []*week, 1)
	go func() {
		retreat := targetDate
		weeks := []*week{m.calculateWeek(targetDate, targetMonth)}
		for {
			retreat = retreat.AddDate(0, 0, -7)
			w := m.calculateWeek(retreat, targetMonth)
			weeks = append([]*week{w}, weeks...)
			if retreat.Month() != targetDate.Month() {
				break
			}
		}
		retWeeks <- weeks
	}()

	advWeeks := make(chan []*week, 1)
	go func() {
		advance := targetDate
		weeks := []*week{}
		for {
			advance = advance.AddDate(0, 0, 7)
			w := m.calculateWeek(advance, targetMonth)
			weeks = append(weeks, w)
			if advance.Month() != targetDate.Month() {
				break
			}
		}
		advWeeks <- weeks
	}()

	ret := <-retWeeks
	adv := <-advWeeks

	m.weeks = append(ret, adv...)
}

func (m *month) calculateWeek(date time.Time, targetMonth time.Month) *week {
	w := &week{}
	done := make(chan struct{}, 2)

	go func() {
		// dateの週の日曜日まで遡る
		retreat := date
		for {
			w.calculateDay(retreat, targetMonth)
			if retreat.Weekday() == time.Sunday {
				break
			}
			retreat = retreat.AddDate(0, 0, -1)
		}
		done <- struct{}{}
	}()

	go func() {
		// dateの週の土曜日まで進む
		advance := date
		for {
			if advance.Weekday() == time.Saturday {
				break
			}
			advance = advance.AddDate(0, 0, 1)
			w.calculateDay(advance, targetMonth)
		}
		done <- struct{}{}
	}()

	<-done
	<-done

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
