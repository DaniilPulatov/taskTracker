package types

import "time"

type Filter struct {
	Day   int
	Month int
	Year  int
}

func NewFilter() *Filter {
	y, m, d := time.Now().Local().Date()
	return &Filter{
		Day:   d,
		Month: int(m),
		Year:  y,
	}
}
