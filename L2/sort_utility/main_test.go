package main

import (
	"math"
	"testing"

	"github.com/merkulovlad/wildberries-L2/sort_utility/internal"
)

func TestNewRecordColumnNumericAndIgnoreTrailing(t *testing.T) {
	opts := internal.Options{
		Column:         2,
		Numeric:        true,
		IgnoreTrailing: true,
	}

	rec := internal.NewRecord("alpha\t123.5   \tbeta", opts)

	if rec.TrimmedKey != "123.5" {
		t.Fatalf("expected TrimmedKey to be %q, got %q", "123.5", rec.TrimmedKey)
	}

	if !rec.NumValid {
		t.Fatalf("expected NumValid to be true")
	}

	if rec.NumericValue != 123.5 {
		t.Fatalf("expected NumericValue to be 123.5, got %v", rec.NumericValue)
	}
}

func TestNewRecordNumericMonthAndHumanParsing(t *testing.T) {
	optsMonth := internal.Options{Month: true}

	validMonth := internal.NewRecord("Feb", optsMonth)
	if !validMonth.MonthValid || validMonth.MonthValue != 2 {
		t.Fatalf("expected month February to be parsed as 2, got (%v, %d)", validMonth.MonthValid, validMonth.MonthValue)
	}

	invalidMonth := internal.NewRecord("Xx", optsMonth)
	if invalidMonth.MonthValid {
		t.Fatalf("expected invalid month string to be rejected")
	}

	optsNumeric := internal.Options{Numeric: true}

	invalidNum := internal.NewRecord("not-a-number", optsNumeric)
	if invalidNum.NumValid {
		t.Fatalf("expected invalid numeric string to be rejected")
	}

	optsHuman := internal.Options{HumanNumeric: true}
	human := internal.NewRecord("1.5M", optsHuman)
	expected := 1.5 * 1024 * 1024

	if !human.HumanNumericValid {
		t.Fatalf("expected valid human numeric string")
	}

	if math.Abs(human.HumanNumericValue-expected) > 1e-9 {
		t.Fatalf("expected human numeric value %.0f, got %v", expected, human.HumanNumericValue)
	}
}

func TestSortRecordsNumericAndReverse(t *testing.T) {
	lines := []string{"10", "2", "001"}

	opts := internal.Options{Numeric: true}
	records := make([]*internal.Record, 0, len(lines))

	for _, line := range lines {
		records = append(records, internal.NewRecord(line, opts))
	}

	internal.SortRecords(records, opts)
	got := []string{records[0].OriginalLine, records[1].OriginalLine, records[2].OriginalLine}
	want := []string{"001", "2", "10"}

	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("expected ascending order %v, got %v", want, got)
		}
	}

	optsReverse := internal.Options{Numeric: true, Reverse: true}
	recordsReverse := make([]*internal.Record, 0, len(lines))

	for _, line := range lines {
		recordsReverse = append(recordsReverse, internal.NewRecord(line, optsReverse))
	}

	internal.SortRecords(recordsReverse, optsReverse)
	got = []string{recordsReverse[0].OriginalLine, recordsReverse[1].OriginalLine, recordsReverse[2].OriginalLine}
	want = []string{"10", "2", "001"}

	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("expected descending order %v, got %v", want, got)
		}
	}
}

func TestCheckIsSorted(t *testing.T) {
	opts := internal.Options{Numeric: true}
	sorted := []*internal.Record{
		internal.NewRecord("1", opts),
		internal.NewRecord("2", opts),
		internal.NewRecord("3", opts),
	}

	if !internal.CheckIsSorted(sorted, opts) {
		t.Fatalf("expected sorted input to be reported as sorted")
	}

	unsorted := []*internal.Record{
		internal.NewRecord("3", opts),
		internal.NewRecord("2", opts),
		internal.NewRecord("1", opts),
	}
	if internal.CheckIsSorted(unsorted, opts) {
		t.Fatalf("expected unsorted input to be reported as unsorted")
	}
}
