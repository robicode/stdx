// These are functions ported from Ruby's Rack::Utils package. Useful utilities
// for HTTP web services.
package httpx

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"time"

	"github.com/robicode/stdx/stringx"
)

// URI escapes. (CGI style space to +)
func Escape(s string) string {
	return url.QueryEscape(s)
}

// Like URI escaping, but with %20 instead of +. Strictly speaking this is
// true URI escaping.
func EscapePath(s string) string {
	return url.PathEscape(s)
}

// Unescapes the **path** component of a URI. See httpx.Unescape() for
// unescaping query parameters or form components.
func UnescapePath(s string) string {
	str, err := url.PathUnescape(s)
	if err != nil {
		return s
	}
	return str
}

// Unescapes a URI escaped string.
// Unlike Rack::Utils, encoding is always UTF-8.
func Unescape(s string) string {
	str, err := url.QueryUnescape(s)
	if err != nil {
		return s
	}
	return str
}

// A QValue represents a quality value header element.
// Used by QValues to return values with their quality
// preferences.
type QValue struct {
	Value   string
	Quality float64
}

type QualityValues []QValue

func (q QualityValues) Len() int {
	return len(q)
}

func (q QualityValues) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
}

func (q QualityValues) Less(i, j int) bool {
	return q[i].Quality < q[j].Quality
}

func (q QValue) String() string {
	return fmt.Sprintf("Value: '%s'; Quality: %f", q.Value, q.Quality)
}

// QValues parses a Q-Value header and returns a set of value-quality
// pairs.
func QValues(header string) []QValue {
	if len(header) == 0 {
		return nil
	}

	parts := stringx.Split(header, regexp.MustCompile(`\s*,\s*`))
	var values []QValue
	var quality float64

	for _, part := range parts {
		var value, params string
		valueParams := stringx.Split(part, regexp.MustCompile(`\s*;\s*`))
		value = valueParams[0]
		if len(valueParams) > 1 {
			params = valueParams[1]
		}
		md := regexp.MustCompile(`\Aq=([\d.]+)`).FindStringSubmatch(params)
		if len(md) == 2 {
			quality, _ = strconv.ParseFloat(md[1], 64)
		} else {
			quality = 1.0
		}
		qv := QValue{
			Value:   value,
			Quality: quality,
		}
		values = append(values, qv)
	}
	return values
}

// HTTPDate formats a time.Time for use in HTTP headers.
func HTTPDate(t time.Time) string {
	return t.Format(http.TimeFormat)
}

type Range struct {
	From int64
	To   int64
}

// Parses the "Range:" header, if present, into an array of Range objects.
// Returns nil if the header is missing or syntactically invalid.
// Returns an empty array if none of the ranges are satisfiable.
func GetByteRanges(rangeHeader string, size int64) []Range {
	// See <http://www.w3.org/Protocols/rfc2616/rfc2616-sec14.html#sec14.35>
	if rangeHeader == "" {
		return nil
	}

	if size <= 0 {
		return []Range{}
	}

	matches := regexp.MustCompile(`bytes=([^;]+)`).FindStringSubmatch(rangeHeader)
	if matches == nil || len(matches) < 2 {
		return nil
	}

	specs := stringx.Split(matches[1], regexp.MustCompile(`,\s*`))

	var (
		ranges []Range
	)

	for _, rangeSpec := range specs {
		rs := regexp.MustCompile(`(\d*)-(\d*)`).FindStringSubmatch(rangeSpec)
		if len(rs) < 3 {
			return nil
		}
		x0, x1 := rs[1], rs[2]
		r0, _ := strconv.Atoi(x0)
		r1, _ := strconv.Atoi(x1)

		if x0 == "" {
			if x1 == "" {
				return nil
			}
			// suffix-byte-range-spec, represents trailing suffix of file
			r0 = int(size - int64(r1))
			if r0 < 0 {
				r0 = 0
			}
			r1 = int(size - 1)
		} else {
			if x1 == "" {
				r1 = int(size - 1)
			} else {
				if r1 < r0 {
					// backwards range is syntactically invalid
					return nil
				}
				if r0 >= int(size) {
					continue
				}
				if r1 >= int(size) {
					r1 = int(size - 1)
				}
			}
		}

		if r0 <= r1 {
			ranges = append(ranges, Range{From: int64(r0), To: int64(r1)})
		}
	}
	return ranges
}
