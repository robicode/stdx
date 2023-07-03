// These are functions ported from Ruby's Rack::Utils package. Useful utilities
// for HTTP web services.
package httpx

import (
	"net/url"
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

// The maximum number of parts a request can contain. Accepting too many part
// can lead to the server running out of file handles.
// Set to `0` for no limit.
const MultipartPartLimit int = 128

func QValues(header string) {

}
