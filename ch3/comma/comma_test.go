package comma

import "testing"

var s [][2]string = [][2]string{
	{"123", "123"},
	{"1234", "1,234"},
	{"20190219", "20,190,219"},
	{"2019219", "2,019,219"},
	{"11223344556677", "11,223,344,556,677"},
	{"13572460910", "13,572,460,910"},
}

func TestComma(t *testing.T) {
	for _, v := range s {
		res := comma(v[0])
		expected := v[1]
		if res != v[1] {
			t.Errorf("res: %q, expected: %q, not same.", res, expected)
		}
	}
}