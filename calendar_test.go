package calma

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var jst = time.FixedZone("JST", +9*60*60)

func TestCalendar_String(t *testing.T) {
	tests := []struct {
		name       string
		date       time.Time
		concurrent bool
		want       string
	}{
		{
			name:       "succeeded",
			date:       time.Date(2021, time.September, 24, 8, 0, 0, 0, jst),
			concurrent: false,
			want: `#### 2021年9月` + "\n" +
				`<font color="red">日</font>|月|火|水|木|金|<font color="blue">土</font>` + "\n" +
				`--------|--------|--------|--------|--------|--------|--------` + "\n" +
				` <font color="red">22</font> | 23 | 24 | 25 | 26 | 27 | <font color="blue">28</font> ` + "\n" +
				` <font color="red">29</font> | 30 | 31 | <b>1 | <b>2 | <b>3 | <font color="blue"><b>4</font> ` + "\n" +
				` <font color="red"><b>5</font> | <b>6 | <b>7 | <b>8 | <b>9 | <b>10 | <font color="blue"><b>11</font> ` + "\n" +
				` <font color="red"><b>12</font> | <b>13 | <b>14 | <b>15 | <b>16 | <b>17 | <font color="blue"><b>18</font> ` + "\n" +
				` <font color="red"><b>19</font> | <font color="red"><b>20</font> | <b>21 | <b>22 | <font color="red"><b>23</font> | <b>24 | <font color="blue"><b>25</font> ` + "\n" +
				` <font color="red"><b>26</font> | <b>27 | <b>28 | <b>29 | <b>30 | 1 | <font color="blue">2</font> ` + "\n",
		},
		{
			name:       "succeeded(concurrently)",
			date:       time.Date(2021, time.September, 24, 8, 0, 0, 0, jst),
			concurrent: true,
			want: `#### 2021年9月` + "\n" +
				`<font color="red">日</font>|月|火|水|木|金|<font color="blue">土</font>` + "\n" +
				`--------|--------|--------|--------|--------|--------|--------` + "\n" +
				` <font color="red">22</font> | 23 | 24 | 25 | 26 | 27 | <font color="blue">28</font> ` + "\n" +
				` <font color="red">29</font> | 30 | 31 | <b>1 | <b>2 | <b>3 | <font color="blue"><b>4</font> ` + "\n" +
				` <font color="red"><b>5</font> | <b>6 | <b>7 | <b>8 | <b>9 | <b>10 | <font color="blue"><b>11</font> ` + "\n" +
				` <font color="red"><b>12</font> | <b>13 | <b>14 | <b>15 | <b>16 | <b>17 | <font color="blue"><b>18</font> ` + "\n" +
				` <font color="red"><b>19</font> | <font color="red"><b>20</font> | <b>21 | <b>22 | <font color="red"><b>23</font> | <b>24 | <font color="blue"><b>25</font> ` + "\n" +
				` <font color="red"><b>26</font> | <b>27 | <b>28 | <b>29 | <b>30 | 1 | <font color="blue">2</font> ` + "\n",
		},
		{
			name: "succeeded(Saturday)",
			date: time.Date(2022, time.June, 18, 8, 0, 0, 0, jst),
			want: `#### 2022年6月` + "\n" +
				`<font color="red">日</font>|月|火|水|木|金|<font color="blue">土</font>` + "\n" +
				`--------|--------|--------|--------|--------|--------|--------` + "\n" +
				` <font color="red">22</font> | 23 | 24 | 25 | 26 | 27 | <font color="blue">28</font> ` + "\n" +
				` <font color="red">29</font> | 30 | 31 | <b>1 | <b>2 | <b>3 | <font color="blue"><b>4</font> ` + "\n" +
				` <font color="red"><b>5</font> | <b>6 | <b>7 | <b>8 | <b>9 | <b>10 | <font color="blue"><b>11</font> ` + "\n" +
				` <font color="red"><b>12</font> | <b>13 | <b>14 | <b>15 | <b>16 | <b>17 | <font color="blue"><b>18</font> ` + "\n" +
				` <font color="red"><b>19</font> | <b>20 | <b>21 | <b>22 | <b>23 | <b>24 | <font color="blue"><b>25</font> ` + "\n" +
				` <font color="red"><b>26</font> | <b>27 | <b>28 | <b>29 | <b>30 | 1 | <font color="blue">2</font> ` + "\n",
		},
		{
			name: "succeeded(Sunday)",
			date: time.Date(2022, time.June, 19, 8, 0, 0, 0, jst),
			want: `#### 2022年6月` + "\n" +
				`<font color="red">日</font>|月|火|水|木|金|<font color="blue">土</font>` + "\n" +
				`--------|--------|--------|--------|--------|--------|--------` + "\n" +
				` <font color="red">29</font> | 30 | 31 | <b>1 | <b>2 | <b>3 | <font color="blue"><b>4</font> ` + "\n" +
				` <font color="red"><b>5</font> | <b>6 | <b>7 | <b>8 | <b>9 | <b>10 | <font color="blue"><b>11</font> ` + "\n" +
				` <font color="red"><b>12</font> | <b>13 | <b>14 | <b>15 | <b>16 | <b>17 | <font color="blue"><b>18</font> ` + "\n" +
				` <font color="red"><b>19</font> | <b>20 | <b>21 | <b>22 | <b>23 | <b>24 | <font color="blue"><b>25</font> ` + "\n" +
				` <font color="red"><b>26</font> | <b>27 | <b>28 | <b>29 | <b>30 | 1 | <font color="blue">2</font> ` + "\n" +
				` <font color="red">3</font> | 4 | 5 | 6 | 7 | 8 | <font color="blue">9</font> ` + "\n",
		},
		{
			// ref: https://www.kantei.go.jp/jp/headline/tokyo2020/shukujitsu.html
			name: "succeeded(July, 1)",
			date: time.Date(2021, time.July, 1, 0, 0, 0, 0, jst),
			want: `#### 2021年7月` + "\n" +
				`<font color="red">日</font>|月|火|水|木|金|<font color="blue">土</font>` + "\n" +
				`--------|--------|--------|--------|--------|--------|--------` + "\n" +
				` <font color="red">20</font> | 21 | 22 | 23 | 24 | 25 | <font color="blue">26</font> ` + "\n" +
				` <font color="red">27</font> | 28 | 29 | 30 | <b>1 | <b>2 | <font color="blue"><b>3</font> ` + "\n" +
				` <font color="red"><b>4</font> | <b>5 | <b>6 | <b>7 | <b>8 | <b>9 | <font color="blue"><b>10</font> ` + "\n" +
				` <font color="red"><b>11</font> | <b>12 | <b>13 | <b>14 | <b>15 | <b>16 | <font color="blue"><b>17</font> ` + "\n" +
				` <font color="red"><b>18</font> | <b>19 | <b>20 | <b>21 | <font color="red"><b>22</font> | <font color="red"><b>23</font> | <font color="blue"><b>24</font> ` + "\n" +
				` <font color="red"><b>25</font> | <b>26 | <b>27 | <b>28 | <b>29 | <b>30 | <font color="blue"><b>31</font> ` + "\n" +
				` <font color="red">1</font> | 2 | 3 | 4 | 5 | 6 | <font color="blue">7</font> ` + "\n",
		},
		{
			// ref: https://www.kantei.go.jp/jp/headline/tokyo2020/shukujitsu.html
			name: "succeeded(October, 31)",
			date: time.Date(2021, time.October, 31, 0, 0, 0, 0, jst),
			want: `#### 2021年10月` + "\n" +
				`<font color="red">日</font>|月|火|水|木|金|<font color="blue">土</font>` + "\n" +
				`--------|--------|--------|--------|--------|--------|--------` + "\n" +
				` <font color="red">26</font> | 27 | 28 | 29 | 30 | <b>1 | <font color="blue"><b>2</font> ` + "\n" +
				` <font color="red"><b>3</font> | <b>4 | <b>5 | <b>6 | <b>7 | <b>8 | <font color="blue"><b>9</font> ` + "\n" +
				` <font color="red"><b>10</font> | <b>11 | <b>12 | <b>13 | <b>14 | <b>15 | <font color="blue"><b>16</font> ` + "\n" +
				` <font color="red"><b>17</font> | <b>18 | <b>19 | <b>20 | <b>21 | <b>22 | <font color="blue"><b>23</font> ` + "\n" +
				` <font color="red"><b>24</font> | <b>25 | <b>26 | <b>27 | <b>28 | <b>29 | <font color="blue"><b>30</font> ` + "\n" +
				` <font color="red"><b>31</font> | 1 | 2 | <font color="red">3</font> | 4 | 5 | <font color="blue">6</font> ` + "\n" +
				` <font color="red">7</font> | 8 | 9 | 10 | 11 | 12 | <font color="blue">13</font> ` + "\n",
		},
	}

	for _, tt := range tests {
		t.Log(tt.name)

		var calendar *Calendar
		var err error
		if tt.concurrent {
			calendar, err = NewCalendarConcurrency(tt.date)
		} else {
			calendar, err = NewCalendar(tt.date)
		}
		assert.NoError(t, err)

		got := fmt.Sprint(calendar)
		assert.Equal(t, tt.want, got)
	}
}

func TestCalendar_HTML(t *testing.T) {
	tests := []struct {
		name string
		date time.Time
		want string
	}{
		{
			name: "succeeded",
			date: time.Date(2022, time.June, 23, 16, 0, 0, 0, jst),
			want: strings.TrimPrefix(`
<h4>2022年6月</h4>

<table>
<thead>
<tr>
<th><font color="red">日</font></th>
<th>月</th>
<th>火</th>
<th>水</th>
<th>木</th>
<th>金</th>
<th><font color="blue">土</font></th>
</tr>
</thead>

<tbody>
<tr>
<td><font color="red">22</font></td>
<td>23</td>
<td>24</td>
<td>25</td>
<td>26</td>
<td>27</td>
<td><font color="blue">28</font></td>
</tr>

<tr>
<td><font color="red">29</font></td>
<td>30</td>
<td>31</td>
<td><b>1</td>
<td><b>2</td>
<td><b>3</td>
<td><font color="blue"><b>4</font></td>
</tr>

<tr>
<td><font color="red"><b>5</font></td>
<td><b>6</td>
<td><b>7</td>
<td><b>8</td>
<td><b>9</td>
<td><b>10</td>
<td><font color="blue"><b>11</font></td>
</tr>

<tr>
<td><font color="red"><b>12</font></td>
<td><b>13</td>
<td><b>14</td>
<td><b>15</td>
<td><b>16</td>
<td><b>17</td>
<td><font color="blue"><b>18</font></td>
</tr>

<tr>
<td><font color="red"><b>19</font></td>
<td><b>20</td>
<td><b>21</td>
<td><b>22</td>
<td><b>23</td>
<td><b>24</td>
<td><font color="blue"><b>25</font></td>
</tr>

<tr>
<td><font color="red"><b>26</font></td>
<td><b>27</td>
<td><b>28</td>
<td><b>29</td>
<td><b>30</td>
<td>1</td>
<td><font color="blue">2</font></td>
</tr>

<tr>
<td><font color="red">3</font></td>
<td>4</td>
<td>5</td>
<td>6</td>
<td>7</td>
<td>8</td>
<td><font color="blue">9</font></td>
</tr>
</tbody>
</table>
`, "\n"),
		},
	}

	for _, tt := range tests {
		t.Log(tt.name)

		calendar, err := NewCalendar(tt.date)
		assert.NoError(t, err)

		got := calendar.HTML()
		assert.Equal(t, tt.want, got)
	}
}
