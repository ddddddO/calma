package main

import (
	"testing"
	"time"
)

func TestBuildCalendar(t *testing.T) {
	tests := []struct {
		name string
		date time.Time
		want string
	}{
		{
			name: "succeeded",
			date: time.Date(2021, time.September, 24, 8, 0, 0, 0, jst),
			want: `#### 2021年9月` + "\n" +
				`|<font color="red">日</font>|月|火|水|木|金|<font color="blue">土</font>|` + "\n" +
				`|--|--|--|--|--|--|--|` + "\n" +
				`| <font color="red">22</font> | 23 | 24 | 25 | 26 | 27 | <font color="blue">28</font> |` + "\n" +
				`| <font color="red">29</font> | 30 | 31 | <b>1 | <b>2 | <b>3 | <font color="blue"><b>4</font> |` + "\n" +
				`| <font color="red"><b>5</font> | <b>6 | <b>7 | <b>8 | <b>9 | <b>10 | <font color="blue"><b>11</font> |` + "\n" +
				`| <font color="red"><b>12</font> | <b>13 | <b>14 | <b>15 | <b>16 | <b>17 | <font color="blue"><b>18</font> |` + "\n" +
				`| <font color="red"><b>19</font> | <font color="red"><b>20</font> | <b>21 | <b>22 | <font color="red"><b>23</font> | <b>24 | <font color="blue"><b>25</font> |` + "\n" +
				`| <font color="red"><b>26</font> | <b>27 | <b>28 | <b>29 | <b>30 | 1 | <font color="blue">2</font> |` + "\n",
		},
	}

	for _, tt := range tests {
		got, err := buildCalendar(tt.date)
		if err != nil {
			t.Error(err)
		}
		if got != tt.want {
			t.Errorf("\ngot: \n%s\nwant: \n%s", got, tt.want)
		}
	}
}
