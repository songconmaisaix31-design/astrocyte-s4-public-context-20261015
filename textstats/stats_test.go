package textstats

import "testing"

func TestUnicodeAndEmpty(t *testing.T) {
	if result := Count("你好 world\n"); result.Words != 2 || result.Characters != 9 || result.Lines != 1 {
		t.Fatalf("unexpected Unicode count: %+v", result)
	}
	if result := Count(""); result != (Stats{}) {
		t.Fatalf("empty input: %+v", result)
	}
	if result := CountAll([]string{"one two", "three"}); result.Words != 3 || result.Lines != 2 {
		t.Fatalf("combined input: %+v", result)
	}
}
