// These are functions ported from Ruby's Rack::Utils package. Useful utilities
// for HTTP web services.
package httpx

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"

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

type QValue struct {
	Value   string
	Quality float64
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
