package stringx

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/huandu/xstrings"
)

type splitType int

const (
	splitTypeAwk splitType = iota
	splitTypeString
	splitTypeRegexp
	splitTypeChars
	splitTypeUnknown
)

// Compares self and other string, ignoring case, and returns
// -1 if other string is larger, 0 if the two are equal, or
// - 1 if other string is smaller.
//
//	CaseCmp("foo", "foo")        // 0
//	CaseCmp("foo", "food")       // -1
//	CaseCmp("food", "foo")       // 1
//	CaseCmp("FOO", "foo")        // 0
//	CaseCmp("foo", "FOO")        // 0
func CaseCmp(str, other string) int {
	if len(str) > len(other) {
		return 1
	} else if len(str) < len(other) {
		return -1
	} else if strings.EqualFold(str, other) {
		return 0
	}

	i := 0
	for ; i < len(str) && i < len(other); i++ {
		if str[i] > other[i] {
			return 1
		} else if str[i] < other[i] {
			return -1
		}
	}
	return 0
}

// Centers str in width.  If width is greater than the length of str,
// returns a new String of length width with str centered and padded with
// pad; otherwise, returns str.
//
//	Center("hello", 4)            // "hello"
//	Center("hello", 20)           // "       hello        "
//	Center("hello", 20, "123")    // "1231231hello12312312"
func Center(s string, width int, pad ...string) string {
	var _pad string
	if pad == nil || len(pad) == 0 {
		_pad = " "
	} else {
		_pad = pad[0]
	}
	return xstrings.Center(s, width, _pad)
}

// Each other_str parameter defines a set of characters to count.  The
// intersection of these sets defines the characters to count in str.  Any
// other_str that starts with a caret ^ is negated.  The sequence c1-c2
// means all characters between c1 and c2.  The backslash character \ can
// be used to escape ^ or - and is otherwise ignored unless it appears at
// the end of a sequence or the end of a other_str.
//
//	str := "hello world"
//	Count(str, "lo")                 // 5
//	Count(str, "lo", "o")            // 2
//	Count(str, "hello", "^l")        // 4
//	Count("hello", "ej-m")           // 4
//
//	str = "hello^world"
//	Count(str, "\\^aeiou")           // 4
//	Count(str, "a\\-eo")             // 3
//
//	str = "hello world\\r\\n"
//	Count(str, "\\")                 // 2
//	Count(str, "\\A")                // 0
//	Count(str, "A-\\w")              // 3
func Count(str string, other_str ...string) int {
	panic("Not yet implemented!")
}

// Returns a copy of str with all characters in the intersection of
// its arguments deleted. Uses the same rules for building the set of
// characters as String#count.
//
//	Delete("hello", "l", "lo")            // "yelow mon"
//	Delete("hello", "lo")                 // "he"
//	Delete("hello", "aeiou", "^e")        // "hell"
//	Delete("hello", "ej-m")               // "ho"
func Delete(s string, pattern ...string) string {
	if pattern == nil || len(pattern) == 0 {
		return s
	}
	switch {
	case len(pattern) == 1:
		return xstrings.Delete(s, pattern[0])
	case len(pattern) == 2:
		delete := xstrings.Delete(s, pattern[0])
		mask := xstrings.Delete(s, pattern[1])
		var bytes []byte
		for _, scp := range s {
			if strings.IndexRune(delete, scp) >= 0 || strings.IndexRune(mask, scp) >= 0 {
				bytes = utf8.AppendRune(bytes, scp)
			}
		}
		return string(bytes)
	}
	return s
}

// Passes each character in str to the given function
func EachRune(str string, block func(r rune)) {
	if block == nil {
		return
	}
	for str != "" {
		r, rsz := utf8.DecodeRuneInString(str)
		block(r)
		str = str[rsz:]
	}
}

// FormatStrings takes a []string and returns a string similar to that
// returned by fmt.Println([]string{...}), but with each element quoted for
// readability.
func FormatStrings(arr []string) string {
	var s string = "["

	for _, item := range arr {
		s = s + fmt.Sprintf("\"%s\", ", item)
	}

	if strings.HasSuffix(s, ", ") {
		s = s[:len(s)-2]
	}
	s = s + "]"

	return s
}

// Insert Inserts the given other string into str; returns the new string.
//
// If the Integer index is positive or zero, inserts other at offset index:
//
//	Insert("foo", 1, "bar")            // "fbaroo"
//
// If the Integer index is negative, counts backward from the end of str
// and inserts other string at offset index+1 (that is, after
// str[index]):
//
//	Insert("foo", -2, "bar")           // "fobaro"
//
// Insert does no utf-8 validation of the other string.
func Insert(str string, index int, other string) string {
	if index == -1 {
		return str + other
	}

	if index >= len(str) {
		return str + other
	}

	var out []byte
	var i int = 0

	if index < 0 {
		index = len(str) + index + 1
	}

	for str != "" {
		if i == index {
			out = append(out, []byte(other)...)
		}
		r, rsz := utf8.DecodeRuneInString(str)
		out = utf8.AppendRune(out, r)
		str = str[rsz:]
		i++
	}

	return string(out)
}

// InsertRune inserts rune r into s at position p.
// No checking is done to validate r is a valid utf-8 rune.
//
// If p is greater than len(s)-1, it is set to len(p)-1.
//
//	str := "helloworld"
//	InsertRune(str, ' ', 5)                 // "hello world"
//	InsertRune(str, '!', 10)                // "helloworld!"
func InsertRune(s string, r rune, p int) string {
	if p > len(s) {
		p = len(s)
	}
	if p < 0 {
		p = len(s) + p + 1
	}
	bytes := []byte(s[:p])
	bytes = utf8.AppendRune(bytes, r)
	return string(bytes) + s[p:]
}

// InsertRunes is like InsertRune, but you can pass multiple runes
// and they will all be inserted at the given position p.
//
// If p is greater than len(s)-1, it is set to len(p)-1.
//
//	str := "helloworld"
//	InsertRunes(str, 5, ',', ' ')                 // "hello, world"
func InsertRunes(s string, p int, rs ...rune) string {
	if rs == nil || len(rs) == 0 {
		return s
	}

	if p > len(s) {
		p = len(s)
	}

	bytes := []byte(s[:p])
	for _, r := range rs {
		bytes = utf8.AppendRune(bytes, r)
	}
	return string(bytes) + s[p:]
}

// IsASCII returns true if s consists entirely of ASCII characters.
// Pulled straight from stdlib and exported.
func IsASCII(s string) bool {
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= utf8.RuneSelf {
			return false
		}
	}
	return true
}

// MatchingStringsets returns true if the two []string slices are equal
// otherwise returns false.
func MatchingStringsets(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// Partition Searches sep or pattern (regexp) in the string and
// returns the part before it, the match, and the part after it. If it is
// not found, returns two empty strings and str.
//
//	Partition("hello", "l")                      // []string{"he", "l", "lo"}
//	Partition("hello", "x")                      // []string{"hello", "", ""}
//	Partition("hello", regexp.MustCompile(`.l`)) // []string{"h", "el", "lo"}
func Partition(str string, pat interface{}) []string {
	if pat == nil {
		return nil
	}

	if sep, ok := pat.(string); ok {
		splits := Split(str, sep, 2)
		if len(splits) < 2 {
			return []string{str, "", ""}
		}
		return []string{splits[0], sep, splits[1]}
	}
	if sep, ok := pat.(*regexp.Regexp); ok {
		match := sep.FindString(str)
		idx := strings.Index(str, match)
		if idx == -1 {
			return []string{str, "", ""}
		}
		return []string{str[:idx], match, str[idx+len(match):]}
	}
	return nil
}

// If the string contains  any invalid byte sequences then replace invalid
// bytes with given replacement string, else returns str. If repl is given
// as a function, replace invalid bytes with returned value of the
// function.
//
//	str := "ab\uFFFDcd\xFF\xCEefg\xFF\xFC\xFD\xFAhijk"
//	Scrub(str)                        // "ab\uFFFDcd\uFFFDefg\uFFFDhijk"
//	Scrub(str, "")                    // "abcdefghijk"
//	Scrub(str, ".")                   // "ab.cd.efg.hijk"
//
//	Scrub(str, func(r rune) string {  // "ab!cd!efg!hijk"
//		return "!"
//	})
func Scrub(str string, repl ...interface{}) string {
	var rep interface{}

	if repl == nil || len(repl) == 0 {
		rep = "\uFFFD"
	} else {
		rep = repl[0]
	}

	if len(repl) > 1 {
		panic("only one scrubber may be defined.")
	}

	if scrubber, ok := rep.(string); ok {
		return xstrings.Scrub(str, scrubber)
	}
	if scrubber, ok := rep.(func(invalid rune) string); ok {
		var buf *strings.Builder
		var r rune
		var size, pos int
		var hasError bool

		origin := str

		for len(str) > 0 {
			r, size = utf8.DecodeRuneInString(str)

			if r == utf8.RuneError {
				if !hasError {
					if buf == nil {
						// buf = &stringBuilder{} // not exposed by xstrings
						buf = &strings.Builder{}
					}

					buf.WriteString(origin[:pos])
					hasError = true
				}
			} else if hasError {
				hasError = false
				buf.WriteString(scrubber(r))

				origin = origin[pos:]
				pos = 0
			}

			pos += size
			str = str[size:]
		}

		if buf != nil {
			buf.WriteString(origin)
			return buf.String()
		}

		// No invalid byte.
		return origin
	}

	panic("invalid scrubber")
}

// Split divides s into substrings based on a pattern, returning a
// slice these substrings.
//
// If pattern is a string, then its contents are used as the
// delimiter when splitting str. If pattern is a single
// space, str is split on whitespace, with leading and trailing
// whitespace and runs of contiguous whitespace characters ignored.
//
// If pattern is a Regexp, str is divided where the
// pattern matches. Whenever the pattern matches a zero-length string,
// str is split into individual characters. If pattern
// contains groups, the respective matches will be returned in the array as
// well.
//
// If pattern is nil, s is split on whitespace as if " " had been passed
// as the pattern.
//
// If the limit parameter is omitted, trailing null fields are
// suppressed. If limit is a positive number, at most that number
// of split substrings will be returned (captured groups will be returned
// as well, but are not counted towards the limit). If limit is
// 1, the entire string is returned as the only entry in an array. If
// negative, there is no limit to the number of fields returned, and
// trailing null fields are not suppressed.
//
// When the input str is empty an empty Array is returned as the string is
// considered to have no fields to split.
//
//	Split(" now's  the time ")                 // []string{"now's", "the", "time"}
//	Split(" now's  the time ", ' ')            // []string{"now's", "the", "time"}
//
//	re := regexp.MustCompile(` `)
//	Split(" now's  the time", re)              // []string{"", "now's", "", "the", "time"}
//
//	re = regexp.MustCompile(`,\s*`)
//	Split("1, 2.34,56, 7", re)                 // []string{"1", "2.34", "56", "7"}
//
//	re = regexp.MustCompile(``)
//	Split("hello", re)                         // []string{"h", "e", "l", "l", "o"}
//	Split("hello", re, 3)                      // []string{"h", "e", "llo"}
//
//	re = regexp.MustCompile(`\s*`)
//	Split("hi mom", re))                       // []string{"h", "i", "m", "o", "m"}
//
//	Split("mellow yellow", "ello")             // []string{"m", "w y", "w"}
//	Split("1,2,,3,4,,", ",")                   // []string{"1", "2", "", "3", "4"}
//	Split("1,2,,3,4,,", ",", 4)                // []string{"1", "2", "", "3,4,,"}
//	Split("1,2,,3,4,,", ",", -4)               // []string{"1", "2", "", "3", "4", "", ""}
//
//	re = regexp.MustCompile(`(:)()()`)
//	Split("1:2:3", re, 2)                      // []string{"1", ":", "", "", "2:3"}
//
//	Split("", ",", -1)                         // []string{}
//
// Currently does not support the block version.
func Split(s string, pattern interface{}, limit ...int) []string {
	if strings.TrimSpace(s) == "" {
		return []string{}
	}

	if pattern == nil {
		pattern = " "
	}

	var lim int
	var emptyCount int = -1

	if len(limit) > 0 {
		lim = limit[0]
	} else {
		lim = -1
	}

	if lim == 1 {
		return []string{s}
	}

	st := getSplitType(pattern, splitTypeRegexp)
	var (
		result []string
	)

	if pat, ok := pattern.(string); ok {
		eptr := utf8.RuneCountInString(s)

		var (
			ptr int = 0
		)

		if limit == nil {
			emptyCount = 0
		}

		if st == splitTypeAwk {
			var (
				end   int  = 0
				skip  bool = true
				start int  = 0
			)

			for ptr < eptr {
				r, rs := utf8.DecodeRuneInString(string(s[ptr]))
				ptr += rs
				if skip {
					if unicode.IsSpace(r) {
						start = ptr
					} else {
						end = ptr - start
						skip = false
					}
				} else if unicode.IsSpace(r) {
					splitString(&result, s, start, end, emptyCount)
					skip = true
					start = ptr
				} else {
					end = ptr - start
				}
				if limit != nil && len(result) >= lim {
					break
				}
			}
			if start != eptr {
				if limit != nil && len(result) == lim {
					ss := result[len(result)-1]
					ss = ss + s[start-1:start+end+1]
					result[len(result)-1] = ss
				} else {
					splitString(&result, s, start, end, emptyCount)
				}
			}
		} else if st == splitTypeChars {
			result = strings.SplitN(s, "", lim)
		} else if st == splitTypeString {
			if limit == nil {
				lim = -1
			}
			result = strings.SplitN(s, pat, lim)
			if limit == nil {
				i := len(result) - 1
				for i > 0 {
					if result[i] == "" {
						result = result[:i]
						i--
						continue
					}
					break
				}
			}
		}
	} else if pat, ok := pattern.(*regexp.Regexp); ok {
		result = reSplit(s, pat, lim)
	}
	return result
}

// Squeeze builds a set of characters from the other_str
// parameter(s) using the procedure described for String#count. Returns a
// new string where runs of the same character that occur in this set are
// replaced by a single character. If no arguments are given, all runs of
// identical characters are replaced by a single character.
//
//	Squeeze("yellow moon", "")             // "yelow mon"
//	Squeeze("  now   is  the", "m-z")      // " now is the"
//	Squeeze("putters shoot balls", " ")    // "puters shot balls"
func Squeeze(s string, other ...string) string {
	var pat string

	if other == nil || len(other) == 0 {
		pat = ""
	} else {
		pat = other[0]
	}

	return xstrings.Squeeze(s, pat)
}

// getSplitType determines the type of split to perform.
func getSplitType(pattern interface{}, defaultType splitType) splitType {
	if pattern == nil {
		return splitTypeAwk
	}
	if pat, ok := pattern.(string); ok {
		if len(pat) == 0 {
			return splitTypeChars
		}
		if len(pat) == 1 && pat[0] == ' ' {
			return splitTypeAwk
		}
		return splitTypeString
	}
	if _, ok := pattern.(*regexp.Regexp); ok {
		return splitTypeRegexp
	}
	return splitTypeUnknown
}

// splitString splits the string str in-place and returns the number of
// empty elements after the call.
func splitString(result *[]string, str string, beg, l, emptyCount int) int {
	if emptyCount >= 0 && l == 0 {
		return emptyCount + 1
	}
	if l > len(str)-1 {
		return -1
	}
	if beg+l > len(str) {
		return -1
	}
	if emptyCount > 0 {
		for emptyCount > 0 {
			if result != nil {
				*result = append(*result, "")
				emptyCount--
			}
		}
	}

	if l < 0 {
		l = len(str) - beg
	}

	if l == 0 {
		*result = append(*result, "")
	} else {
		s := str[beg : beg+l]
		*result = append(*result, s)
	}

	return emptyCount
}

// reSplit splits a string based on a *regexp.Regexp, accounting for
// capture groups.
func reSplit(s string, re *regexp.Regexp, n int) []string {
	if n == 0 {
		return nil
	}

	if len(re.String()) > 0 && len(s) == 0 {
		return []string{""}
	}

	matches := re.FindAllStringIndex(s, n)
	strings := make([]string, 0, len(matches))

	beg := 0
	end := 0
	for _, match := range matches {
		if n > 0 && len(strings) >= n-1 {
			break
		}

		end = match[0]
		if match[1] != 0 {
			strings = append(strings, s[beg:end])
		}
		beg = match[1]
	}

	subs := re.FindStringSubmatch(s)
	if subs != nil {
		for i := 1; i < len(subs); i++ {
			strings = append(strings, subs[i])
		}
	}

	if end != len(s) {
		strings = append(strings, s[beg:])
	}

	return strings
}

// Decodes str (which may contain binary data) according to the
// format string, returning an array of each value extracted. The format
// string consists of a sequence of single-character directives, summarized
// in the table at the end of this entry. Each directive may be followed by
// a number, indicating the number of times to repeat with this directive.
// An asterisk (“*”) will use up all remaining elements. The directives
// sSiIlL may each be followed by an underscore (“_”) or exclamation mark
// (“!”) to use the underlying platform's native size for the specified
// type; otherwise, it uses a platform-independent consistent size. Spaces
// are ignored in the format string.
//
// Note that unlike the Ruby version of this function, this function always returns
// a slice of []string, because as far as I know Go does not support variable types
// in the same slice.
//
//	Unpack("abc \0\0abc \0\0", "A6Z6")          // []string{"abc", "abc "}
//	Unpack("abc \0\0", "a3a3")                  // []string{"abc", " \000\000"}
//	Unpack("abc \0abc \0", "Z*Z*")              // []string{"abc ", "abc "}
//	Unpack("aa", "b8B8")                        // []string{"10000110", "01100001"}
//	Unpack("aaa", "h2H2c")                      // []string{"16", "61", "97"}
//	Unpack("\xfe\xff\xfe\xff", "sS")            // []string{"-2", "65534"}
//	Unpack("now=20is", "M*")                    // []string{"now is"}
//	Unpack("whole", "xax2aX2aX1aX2a")           // []string{"h", "e", "l", "l", "o"}
//
// The README for this package summarizes the various formats. Note that unlike the Ruby
// version of this function, we always return a slice of []string, because as far as I know
// Go does not support multiple variable types in the same slice.
func Unpack(str, format string) []string {
	return nil
}
