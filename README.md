# <font color="red">cal</font><font color="blue">ma</font>
Calendar for Markdown<br>
Inspired by [jpcal](https://github.com/y-yagi/jpcal)<br>
&#x26a0;Only for Japan&#x26a0;

## Usage

```console
19:07:38 > calma -h
This CLI outputs Japanese calendar in Markdown. It supports national holidays.

Usage of calma:
  -a int
        Number of months to advance
  -r int
        Number of months to retreat

19:07:57 > calma
#### 2021年9月
|<font color="red">日</font>|月|火|水|木|金|<font color="blue">土</font>|
|--|--|--|--|--|--|--|
| <font color="red">22</font> | 23 | 24 | 25 | 26 | 27 | <font color="blue">28</font> |
| <font color="red">29</font> | 30 | 31 | <b>1 | <b>2 | <b>3 | <font color="blue"><b>4</font> |
| <font color="red"><b>5</font> | <b>6 | <b>7 | <b>8 | <b>9 | <b>10 | <font color="blue"><b>11</font> |
| <font color="red"><b>12</font> | <b>13 | <b>14 | <b>15 | <b>16 | <b>17 | <font color="blue"><b>18</font> |
| <font color="red"><b>19</font> | <font color="red"><b>20</font> | <b>21 | <b>22 | <font color="red"><b>23</font> | <b>24 | <font color="blue"><b>25</font> |
| <font color="red"><b>26</font> | <b>27 | <b>28 | <b>29 | <b>30 | 1 | <font color="blue">2</font> |
```

Markdown output image<br>
![image](https://github.com/ddddddO/calma/blob/main/sample.png)

## Installation

go version is 1.16 or higher.

```console
go install github.com/ddddddO/calma@latest
```

go version is 1.15 or less.
```console
go get github.com/ddddddO/calma
```

or using Homebrew.
```console
brew install ddddddO/tap/calma
```

**or, download binary from [here](https://github.com/ddddddO/calma/releases).**
