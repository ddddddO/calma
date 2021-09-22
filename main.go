package main

import (
	"fmt"
	"os"
	"text/template"
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
	wk := week{Sunday: 1, Monday: 2, Tuesday: 3, Wednesday: 4, Thursday: 5, Friday: 6, Saturday: 7}
	tmpl, err := template.New("calma").Parse(calendarTemplate)
	if err != nil {
		fmt.Errorf("%+v", err)
		os.Exit(1)
	}
	if err := tmpl.Execute(os.Stdout, wk); err != nil {
		fmt.Errorf("%+v", err)
		os.Exit(1)
	}

}

const calendarTemplate = `
|<font color="red">日</font>|月|火|水|木|金|<font color="blue">土</font>|
|--|--|--|--|--|--|--|
|{{.Sunday}}|{{.Monday}}|{{.Tuesday}}|{{.Wednesday}}|{{.Thursday}}|{{.Friday}}|{{.Saturday}}|
`
