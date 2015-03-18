package naturalsort

import (
	"reflect"
	"sort"
	"testing"
)

func TestSortValid(t *testing.T) {
	cases := []struct {
		data, expected []string
	}{
		{
			nil,
			nil,
		},
		{
			[]string{},
			[]string{},
		},
		{
			[]string{"a"},
			[]string{"a"},
		},
		{
			[]string{"0"},
			[]string{"0"},
		},
		{
			[]string{"1", "2", "30", "22", "0", "00", "3"},
			[]string{"0", "00", "1", "2", "3", "22", "30"},
		},
		{
			[]string{"A1", "A0", "A21", "A11", "A111", "A2"},
			[]string{"A0", "A1", "A2", "A11", "A21", "A111"},
		},
		{
			[]string{"A1BA1", "A11AA1", "A2AB0", "B1AA1", "A1AA1"},
			[]string{"A1AA1", "A1BA1", "A2AB0", "A11AA1", "B1AA1"},
		},
	}

	for i, c := range cases {
		sort.Sort(NaturalSort(c.data))
		if !reflect.DeepEqual(c.data, c.expected) {
			t.Fatalf("Wrong order in test case #%d.\nExpected=%v\nGot=%v", i, c.expected, c.data)
		}
	}

}

func BenchmarkSort(b *testing.B) {
	var data = [...]string{"A1BA1", "A11AA1", "A2AB0", "B1AA1", "A1AA1"}
	for ii := 0; ii < b.N; ii++ {
		d := NaturalSort(data[:])
		sort.Sort(d)
	}
}
