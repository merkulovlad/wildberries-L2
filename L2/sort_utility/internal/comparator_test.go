package internal

import "testing"

func TestSortRecordsMonthOrdering(t *testing.T) {
	opts := Options{Month: true}
	lines := []string{"Mar", "zzz", "Jan"}

	records := make([]*Record, 0, len(lines))
	for _, line := range lines {
		records = append(records, NewRecord(line, opts))
	}

	SortRecords(records, opts)
	got := []string{records[0].OriginalLine, records[1].OriginalLine, records[2].OriginalLine}
	want := []string{"Jan", "Mar", "zzz"}

	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("month sort order mismatch: got %v want %v", got, want)
		}
	}
}

func TestNumericLessPrefersValidNumbers(t *testing.T) {
	opts := Options{Numeric: true}
	valid := NewRecord("10", opts)
	invalid := NewRecord("abc", opts)

	if !less(valid, invalid, opts) {
		t.Fatalf("expected valid numeric record to be considered less; valid.NumValid=%v invalid.NumValid=%v", valid.NumValid, invalid.NumValid)
	}

	if less(invalid, valid, opts) {
		t.Fatalf("expected invalid numeric record not to precede valid ones")
	}
}

func TestCheckIsSortedReverseNumeric(t *testing.T) {
	opts := Options{Numeric: true, Reverse: true}
	desc := []*Record{
		NewRecord("9", opts),
		NewRecord("3", opts),
		NewRecord("1", opts),
	}

	if !CheckIsSorted(desc, opts) {
		t.Fatalf("expected descending numeric slice to be sorted with reverse option")
	}

	asc := []*Record{
		NewRecord("1", opts),
		NewRecord("3", opts),
		NewRecord("9", opts),
	}
	if CheckIsSorted(asc, opts) {
		t.Fatalf("expected ascending numeric slice to be unsorted with reverse option")
	}
}

func TestHumanNumericComparison(t *testing.T) {
	opts := Options{HumanNumeric: true}
	lines := []string{"512K", "1M", "32K"}
	records := make([]*Record, 0, len(lines))

	for _, line := range lines {
		records = append(records, NewRecord(line, opts))
	}

	SortRecords(records, opts)
	got := []string{records[0].OriginalLine, records[1].OriginalLine, records[2].OriginalLine}
	want := []string{"32K", "512K", "1M"}

	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("human numeric sort order mismatch: got %v want %v", got, want)
		}
	}
}
