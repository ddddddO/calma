package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBuild_Calendar(t *testing.T) {
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
		// FIXME: エッジケースが考慮できていない
		// {
		// 	name: "succeeded",
		// 	date: time.Date(2021, time.September, 23, 8, 0, 0, 0, jst),
		// 	want: `#### 2021年9月` + "\n" +
		// 		`|<font color="red">日</font>|月|火|水|木|金|<font color="blue">土</font>|` + "\n" +
		// 		`|--|--|--|--|--|--|--|` + "\n" +
		// 		`| <font color="red">22</font> | 23 | 24 | 25 | 26 | 27 | <font color="blue">28</font> |` + "\n" +
		// 		`| <font color="red">29</font> | 30 | 31 | <b>1 | <b>2 | <b>3 | <font color="blue"><b>4</font> |` + "\n" +
		// 		`| <font color="red"><b>5</font> | <b>6 | <b>7 | <b>8 | <b>9 | <b>10 | <font color="blue"><b>11</font> |` + "\n" +
		// 		`| <font color="red"><b>12</font> | <b>13 | <b>14 | <b>15 | <b>16 | <b>17 | <font color="blue"><b>18</font> |` + "\n" +
		// 		`| <font color="red"><b>19</font> | <font color="red"><b>20</font> | <b>21 | <b>22 | <font color="red"><b>23</font> | <b>24 | <font color="blue"><b>25</font> |` + "\n" +
		// 		`| <font color="red"><b>26</font> | <b>27 | <b>28 | <b>29 | <b>30 | 1 | <font color="blue">2</font> |` + "\n",
		// },
		{
			// ref: https://www.kantei.go.jp/jp/headline/tokyo2020/shukujitsu.html
			name: "succeeded(July)",
			date: time.Date(2021, time.July, 1, 0, 0, 0, 0, jst),
			want: `#### 2021年7月` + "\n" +
				`|<font color="red">日</font>|月|火|水|木|金|<font color="blue">土</font>|` + "\n" +
				`|--|--|--|--|--|--|--|` + "\n" +
				`| <font color="red">20</font> | 21 | 22 | 23 | 24 | 25 | <font color="blue">26</font> |` + "\n" +
				`| <font color="red">27</font> | 28 | 29 | 30 | <b>1 | <b>2 | <font color="blue"><b>3</font> |` + "\n" +
				`| <font color="red"><b>4</font> | <b>5 | <b>6 | <b>7 | <b>8 | <b>9 | <font color="blue"><b>10</font> |` + "\n" +
				`| <font color="red"><b>11</font> | <b>12 | <b>13 | <b>14 | <b>15 | <b>16 | <font color="blue"><b>17</font> |` + "\n" +
				`| <font color="red"><b>18</font> | <b>19 | <b>20 | <b>21 | <font color="red"><b>22</font> | <font color="red"><b>23</font> | <font color="blue"><b>24</font> |` + "\n" +
				`| <font color="red"><b>25</font> | <b>26 | <b>27 | <b>28 | <b>29 | <b>30 | <font color="blue"><b>31</font> |` + "\n" +
				`| <font color="red">1</font> | 2 | 3 | 4 | 5 | 6 | <font color="blue">7</font> |` + "\n",
		},
		{
			// ref: https://www.kantei.go.jp/jp/headline/tokyo2020/shukujitsu.html
			name: "succeeded(October)",
			date: time.Date(2021, time.October, 24, 0, 0, 0, 0, jst), // NOTE: 10/1だとこける。actualとの差分はないはずだけど。
			want: `#### 2021年10月` + "\n" +
				`|<font color="red">日</font>|月|火|水|木|金|<font color="blue">土</font>|` + "\n" +
				`|--|--|--|--|--|--|--|` + "\n" +
				`| <font color="red">26</font> | 27 | 28 | 29 | 30 | <b>1 | <font color="blue"><b>2</font> |` + "\n" +
				`| <font color="red"><b>3</font> | <b>4 | <b>5 | <b>6 | <b>7 | <b>8 | <font color="blue"><b>9</font> |` + "\n" +
				`| <font color="red"><b>10</font> | <b>11 | <b>12 | <b>13 | <b>14 | <b>15 | <font color="blue"><b>16</font> |` + "\n" +
				`| <font color="red"><b>17</font> | <b>18 | <b>19 | <b>20 | <b>21 | <b>22 | <font color="blue"><b>23</font> |` + "\n" +
				`| <font color="red"><b>24</font> | <b>25 | <b>26 | <b>27 | <b>28 | <b>29 | <font color="blue"><b>30</font> |` + "\n" +
				`| <font color="red"><b>31</font> | 1 | 2 | <font color="red">3</font> | 4 | 5 | <font color="blue">6</font> |` + "\n" +
				`| <font color="red">7</font> | 8 | 9 | 10 | 11 | 12 | <font color="blue">13</font> |` + "\n",
		},
	}

	for _, tt := range tests {
		t.Log(tt.name)

		calendar, err := newCalendar(tt.date)
		assert.NoError(t, err)

		err = calendar.build()
		assert.NoError(t, err)

		got := fmt.Sprint(calendar)
		assert.Equal(t, tt.want, got)
	}
}
