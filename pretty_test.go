package pp_test

import (
	"io"
	"testing"
	"time"

	"github.com/nochso/pp"
)

var tests = []struct {
	Input  interface{}
	Output string
}{
	{
		Input:  0,
		Output: "0",
	},
	{
		Input: []int{0, 1, 2},
		Output: `[]int{
	0,
	1,
	2,
}`,
	},
	{
		Input:  "x",
		Output: `"x"`,
	},
	{
		Input: map[string]string{"foo": "bar"},
		Output: `map[string]string{
	"foo": "bar",
}`,
	},
	{
		Input:  time.Date(2022, time.February, 13, 11, 2, 52, 0, time.UTC),
		Output: "time.Date(2022, time.February, 13, 11, 2, 52, 0, time.UTC)",
	},
	{
		Input:  io.Discard,
		Output: "io.discard{}",
	},
	{
		Input:  true,
		Output: "true",
	},
	{
		Input:  byte(5),
		Output: "0x5",
	},
	{
		Input: []byte("abc"),
		Output: `[]byte{
	0x61,
	0x62,
	0x63,
}`,
	},
	{
		Input: struct {
			name  string
			score float32
		}{
			name:  "foo",
			score: 22 / 7,
		},
		Output: `struct {
	name  string
	score float32
}{
	name:  "foo",
	score: 3,
}`,
	},
}

func TestSprint(t *testing.T) {
	for i, test := range tests {
		actual := pp.Sprint(test.Input)
		if actual != test.Output {
			t.Errorf("Test %d input: %#v\nexpected output:\n'%s'\nactual:\n'%s'", i, test.Input, test.Output, actual)
		}
	}
}
