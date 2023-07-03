package httpx

import (
	"testing"
)

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

func checkQValues(t *testing.T, expected, actual []QValue) []qvDiff {
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
