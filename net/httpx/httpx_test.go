package httpx

import (
	"fmt"
	"strings"
	"testing"
)

func Test_GetByteRanges(t *testing.T) {
	// Missing or invalid byte ranges
	checkByteRanges(t, "ignore missing or syntactically invalid byte ranges", nil, GetByteRanges("", 500))
	checkByteRanges(t, "ignore missing or syntactically invalid byte ranges", nil, GetByteRanges("foobar", 500))
	checkByteRanges(t, "ignore missing or syntactically invalid byte ranges", nil, GetByteRanges("furlongs=123-456", 500))
	checkByteRanges(t, "ignore missing or syntactically invalid byte ranges", nil, GetByteRanges("bytes=", 500))
	checkByteRanges(t, "ignore missing or syntactically invalid byte ranges", nil, GetByteRanges("bytes=-", 500))
	checkByteRanges(t, "ignore missing or syntactically invalid byte ranges", nil, GetByteRanges("bytes=123,456", 500))

	// A range of non-positive length is syntactically invalid and ignored:
	checkByteRanges(t, "ignore missing or syntactically invalid byte ranges", nil, GetByteRanges("bytes=456-123", 500))
	checkByteRanges(t, "ignore missing or syntactically invalid byte ranges", nil, GetByteRanges("bytes=456-455", 500))

	// Simple byte ranges
	checkByteRanges(t, "parses simple byte ranges (1)", []Range{{From: 123, To: 456}}, GetByteRanges("bytes=123-456", 500))
	checkByteRanges(t, "parses simple byte ranges (2)", []Range{{From: 123, To: 499}}, GetByteRanges("bytes=123-", 500))
	checkByteRanges(t, "parses simple byte ranges (3)", []Range{{From: 400, To: 499}}, GetByteRanges("bytes=-100", 500))
	checkByteRanges(t, "parses simple byte ranges (4)", []Range{{From: 0, To: 0}}, GetByteRanges("bytes=0-0", 500))
	checkByteRanges(t, "parses simple byte ranges (5)", []Range{{From: 499, To: 499}}, GetByteRanges("bytes=499-499", 500))

	// Multiple byte ranges
	checkByteRanges(t, "parses several byte ranges", []Range{{From: 500, To: 600}, {From: 601, To: 999}}, GetByteRanges("bytes=500-600,601-999", 1000))

	// Truncation of byte ranges
	checkByteRanges(t, "truncates byte ranges (1)", []Range{{From: 123, To: 499}}, GetByteRanges("bytes=123-999", 500))
	checkByteRanges(t, "truncates byte ranges (2)", []Range{}, GetByteRanges("bytes=600-999", 500))
	checkByteRanges(t, "truncates byte ranges (3)", []Range{{From: 0, To: 499}}, GetByteRanges("bytes=-999", 500))

	// Unsatisfiable byte ranges
	checkByteRanges(t, "ignores unsatisfiable byte ranges (1)", []Range{}, GetByteRanges("bytes=500-501", 500))
	checkByteRanges(t, "ignores unsatisfiable byte ranges (2)", []Range{}, GetByteRanges("bytes=500-", 500))
	checkByteRanges(t, "ignores unsatisfiable byte ranges (3)", []Range{}, GetByteRanges("bytes=999-", 500))
	checkByteRanges(t, "ignores unsatisfiable byte ranges (4)", []Range{}, GetByteRanges("bytes=-0", 500))

	// handle byte ranges of empty files
	checkByteRanges(t, "handle byte ranges of empty files (1)", []Range{}, GetByteRanges("bytes=123-456", 0))
	checkByteRanges(t, "handle byte ranges of empty files (2)", []Range{}, GetByteRanges("bytes=0-", 0))
	checkByteRanges(t, "handle byte ranges of empty files (3)", []Range{}, GetByteRanges("bytes=-100", 0))
	checkByteRanges(t, "handle byte ranges of empty files (4)", []Range{}, GetByteRanges("bytes=0-0", 0))
	checkByteRanges(t, "handle byte ranges of empty files (5)", []Range{}, GetByteRanges("bytes=-0", 0))
}

func Test_QValues(t *testing.T) {
	h := `text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8`
	expected := []QValue{
		{Value: "text/html", Quality: 1},
		{Value: "application/xhtml+xml", Quality: 1},
		{Value: "application/xml", Quality: 0.9},
		{Value: "image/avif", Quality: 1},
		{Value: "image/webp", Quality: 1},
		{Value: "*/*", Quality: 0.8},
	}

	values := QValues(h)

	checkQValues(t, expected, values)
}

type qvDiff struct {
	Expected QValue
	Actual   QValue
}

// checkQValues checks the return of QValues against the expected results and fails if there
// is a mismatch.
func checkQValues(t *testing.T, expected, actual []QValue) []qvDiff {
	if len(expected) != len(actual) {
		t.Errorf("expected length of expected (%d) to match actual (%d) values.", len(expected), len(actual))
		return []qvDiff{}
	}
	var errors []qvDiff
	for i, value := range expected {
		if value.Value != actual[i].Value {
			errors = append(errors, qvDiff{
				Expected: value,
				Actual:   actual[i],
			})
			t.Errorf("expected and actual values differ:\nActual: %v\nExpected: %v\n", actual, expected)
		}
	}
	return errors
}

// checkByteRanges checks the return of a GetByteRanges call against expected return
// values and fails the test if there is a mismatch.
func checkByteRanges(t *testing.T, description string, expected, actual []Range) {
	var errors []string
	if expected == nil && actual != nil {
		t.Errorf("expected nil but got %v", actual)
		return
	}

	if len(expected) != len(actual) {
		fmt.Printf("Test '%s':\n", strings.Title(description))
		t.Errorf("length of expected (%d) and actual (%d) results do not match\nExpected:%v\nActual:%v\n", len(expected), len(actual), expected, actual)
		return
	}

	for i, byterange := range expected {
		if actual[i].From != byterange.From {
			errors = append(errors, fmt.Sprintf("expected byte range %d From value to be %d but was %d", i, byterange.From, actual[i].From))
		}
		if actual[i].To != byterange.To {
			errors = append(errors, fmt.Sprintf("expected byte range %d To value to be %d but was %d", i, byterange.To, actual[i].To))
		}
	}

	if len(errors) != 0 {
		fmt.Printf("Test: '%s' Errors (%d):\n", strings.Title(description), len(errors))
	}
	for _, err := range errors {
		t.Error(err)
	}
}
