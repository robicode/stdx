// These are functions ported from Ruby's Rack::Utils package. Useful utilities
// for HTTP web services.
package httpx

import (
	"log"
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

// The maximum number of parts a request can contain. Accepting too many part
// can lead to the server running out of file handles.
// Set to `0` for no limit.
const MultipartPartLimit int = 128

// QValues parses a Q-Value header and returns a set of value-quality
// pairs.
func QValues(header string) []QValue {
	if len(header) == 0 {
		return nil
	}
	parts := stringx.Split(header, regexp.MustCompile(`\s*,\s*`))
	var values []QValue

	for _, part := range parts {
		valueParams := stringx.Split(part, `\s*;\s*`, 2)
		if valueParams == nil || len(valueParams) < 2 {
			continue
		}
		value := valueParams[0]
		parameters := valueParams[1]
		quality := 1.0
		md := regexp.MustCompile(`\Aq=([\d.]+)`).FindStringSubmatch(parameters)
		if md == nil || len(md) != 1 {
			continue
		}
		quality, err := strconv.ParseFloat(md[0], 64)
		if err != nil {
			log.Printf("quality is not a valid floating point number: %v", md[1])
			continue
		}
		qv := QValue{
			Value:   value,
			Quality: quality,
		}
		values = append(values, qv)
	}
	return values
}
