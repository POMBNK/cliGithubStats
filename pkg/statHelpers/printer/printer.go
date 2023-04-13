package printer

import (
	"fmt"
	"sort"
	"time"
)

const (
	daysInLastSixMonths  = 183
	weeksInLastSixMonths = 26
)

type column []int

type Printer struct {
}

func New() *Printer {
	return &Printer{}
}

func (p *Printer) Show(commits map[int]int) {
	keys := p.sortMapIntoSlice(commits)
	cols := p.buildCols(keys, commits)
	p.printCells(cols)
}

func (p *Printer) sortMapIntoSlice(m map[int]int) []int {
	// order map
	// To store the keys in slice in sorted order
	var keys []int
	for k := range m {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	return keys
}

func (p *Printer) buildCols(keys []int, commits map[int]int) map[int]column {
	cols := make(map[int]column)
	col := column{}

	for _, k := range keys {
		week := int(k / 7) //26,25...1
		dayinweek := k % 7 // 0,1,2,3,4,5,6

		if dayinweek == 0 { //reset
			col = column{}
		}

		col = append(col, commits[k])

		if dayinweek == 6 {
			cols[week] = col
		}
	}

	return cols
}

func (p *Printer) printCells(cols map[int]column) {
	p.printMonths()
	for j := 6; j >= 0; j-- {
		for i := weeksInLastSixMonths + 1; i > 0; i-- { //>= ?
			if i == weeksInLastSixMonths+1 {
				p.printDayCol(j)
			}
			if col, ok := cols[i]; ok {
				//special case today
				if i == 0 && j == p.calcOffset()-1 {
					p.printCell(col[j], true)
					continue
				} else {
					if len(col) > j {
						p.printCell(col[j], false)
						continue
					}
				}
			}
			p.printCell(0, false)
		}
		fmt.Printf("\n")
	}
}

func (p *Printer) printCell(val int, today bool) {
	escape := "\033[0;37;30m"
	switch {
	case val > 0 && val < 5:
		escape = "\033[1;30;47m"
	case val >= 5 && val < 10:
		escape = "\033[1;30;43m"
	case val >= 10:
		escape = "\033[1;30;42m"
	}

	if today {
		escape = "\033[1;37;45m"
	}

	if val == 0 {
		fmt.Printf(escape + "  - " + "\033[0m")
		return
	}

	str := "  %d "
	switch {
	case val >= 10:
		str = " %d "
	case val >= 100:
		str = "%d "
	}

	fmt.Printf(escape+str+"\033[0m", val)
}

func (p *Printer) printMonths() {
	week := p.getBeginningOfDay(time.Now()).Add(-(daysInLastSixMonths * time.Hour * 24))
	month := week.Month()
	fmt.Printf("         ")
	for {
		if week.Month() != month {
			fmt.Printf("%s ", week.Month().String()[:3])
			month = week.Month()
		} else {
			fmt.Printf("    ")
		}

		week = week.Add(7 * time.Hour * 24)
		if week.After(time.Now()) {
			break
		}
	}
	fmt.Printf("\n")
}

func (p *Printer) getBeginningOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	startOfDay := time.Date(year, month, day, 0, 0, 0, 0, t.Location())
	return startOfDay
}

func (p *Printer) printDayCol(day int) {
	out := "     "
	switch day {
	case 1:
		out = " Mon "
	case 3:
		out = " Wed "
	case 5:
		out = " Fri "
	}

	fmt.Printf(out)
}
func (p *Printer) calcOffset() int {
	var offset int
	weekday := time.Now().Weekday()

	switch weekday {

	case time.Monday:
		offset = 7
	case time.Tuesday:
		offset = 6
	case time.Wednesday:
		offset = 5
	case time.Thursday:
		offset = 4
	case time.Friday:
		offset = 3
	case time.Saturday:
		offset = 2
	case time.Sunday:
		offset = 1
	}

	return offset
}
