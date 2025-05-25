package main

import (
	"strconv"
	"syscall/js" // nolint
	"text/template"
	"time"

	"github.com/ddddddO/calma"
)

func main() {
	c := make(chan struct{}, 0)
	println("calma WebAssembly Initialized")
	registerCallbacks()
	<-c
}

func registerCallbacks() {
	js.Global().Set("generateCalender", js.FuncOf(generateCalender))
}

func generateCalender(this js.Value, args []js.Value) interface{} {
	document := js.Global().Get("document")
	getElementByID := getElementByIDFunc(document)

	y := getElementByID("year").Get("value").String()
	m := getElementByID("month").Get("value").String()
	year, err := strconv.Atoi(y)
	if err != nil {
		alert(err.Error())
		return nil
	}
	month, err := strconv.Atoi(m)
	if err != nil {
		alert(err.Error())
		return nil
	}

	calender, err := calma.NewCalendar(targetTime(year, month))
	if err != nil {
		alert(err.Error())
		return nil
	}

	div := getElementByID("result")
	if prePre := getElementByID("redered_calender"); !prePre.IsNull() {
		removeChildFunc(div)(prePre)
	}

	pre := createElementFunc(document)("pre")
	pre.Set("id", "redered_calender")
	pre.Set("innerHTML", template.HTMLEscapeString(calender.String()))
	appendChildFunc(div)(pre)
	appendChildFunc(getElementByID("main"))(div)

	return nil
}

const middleDay = 14

func targetTime(year int, month int) time.Time {
	return time.Date(year, time.Month(month), middleDay, 0, 0, 0, 0, time.Local)
}

func getElementByIDFunc(document js.Value) func(id string) js.Value {
	return func(id string) js.Value {
		return document.Call("getElementById", id)
	}
}

func createElementFunc(document js.Value) func(element string) js.Value {
	return func(element string) js.Value {
		return document.Call("createElement", element)
	}
}

func removeChildFunc(element js.Value) func(target js.Value) {
	return func(target js.Value) {
		element.Call("removeChild", target)
	}
}

func appendChildFunc(element js.Value) func(target js.Value) {
	return func(target js.Value) {
		element.Call("appendChild", target)
	}
}

func alert(msg string) {
	js.Global().Call("alert", msg)
}
