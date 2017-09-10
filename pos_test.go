package main

import "testing"

func TestDir(t *testing.T) {
	type tableTest struct {
		pos position
		dir direction
	}
	table := []tableTest{
		{position{3, 2}, E},
		{position{4, 1}, ENE},
		{position{3, 1}, NE},
		{position{3, 0}, NNE},
		{position{2, 1}, N},
		{position{1, 0}, NNW},
		{position{1, 1}, NW},
		{position{0, 1}, WNW},
		{position{1, 2}, W},
		{position{0, 3}, WSW},
		{position{1, 3}, SW},
		{position{1, 4}, SSW},
		{position{2, 3}, S},
		{position{3, 4}, SSE},
		{position{3, 3}, SE},
		{position{4, 3}, ESE},
	}
	for _, test := range table {
		if test.pos.Dir(position{2, 2}) != test.dir {
			t.Errorf("Bad direction for %+v\n", test)
		}
	}
}
