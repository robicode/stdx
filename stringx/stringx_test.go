package stringx

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
)

var errors int = 0

func Test_Center(t *testing.T) {
	if Center("hello", 4) != "hello" {
		t.Error("expected Center() to return the original string")
	}
	if Center("hello", 20) != "       hello        " {
		t.Errorf("expected '       hello        ' but got '%s'", Center("hello", 20))
	}
	if Center("hello", 20, "123") != "1231231hello12312312" {
		t.Errorf("expected '1231231hello12312312' but got '%s'", Center("hello", 20, "123"))
	}
}

// func Test_Count(t *testing.T) {
// TODO implement a Ruby-style Count() function
// Note: Commented tests are not passing yet.
// str := "hello world"
// if Count(str, "lo") != 5 {
// 	t.Errorf("expected Count() to be 5 but was %d", Count(str, "lo"))
// }

// xstrings does not support multiple patterns, as per the following two (failing) tests
// and I haven't implemented a fix yet:
// if Count(str, "lo", "o") != 2 {
// 	t.Errorf("expected Count() to be 2 but was %d", Count(str, "lo", "o"))
// }
// if Count(str, "hello", "^l") != 4 {
// 	t.Errorf("expected Count() to be 4 but was %d", Count(str, "hello", "^l"))
// }

// if Count(str, "ej-m") != 4 {
// 	t.Errorf("expected Count() to be %d but was %d", 4, Count(str, "ej-m"))
// }

// str = "hello^world"
// if Count(str, "\\^aeiou") != 4 {
// 	t.Errorf("expected Count() to be %d but was %d", 4, Count(str, "\\^aeiou"))
// }
// if Count(str, "a\\-eo") != 3 {
// 	t.Errorf("expected Count() to be %d but was %d", 3, Count(str, "a\\-eo"))
// }

// str = "hello world\\r\\n"
// if Count(str, "\\") != 2 {
// 	t.Errorf("expected Count('%s', '\\\\') to be %d but was %d", str, 2, Count(str, "\\"))
// }
// if Count(str, "\\A") != 0 {
// 	t.Errorf("expected Count('%s', '\\\\A') to be %d but was %d", str, 0, Count(str, "\\A"))
// }
// if Count(str, "A-\\w") != 3 {
// 	t.Errorf("expected Count('%s', 'A-\\\\w') to be %d but was %d", str, 3, Count(str, "A-\\w"))
// }
// }

func Test_Delete(t *testing.T) {
	if Delete("hello", "l", "lo") != "heo" {
		t.Errorf("expected 'heo' but got '%s'", Delete("hello", "l", "lo"))
	}
	if Delete("hello", "lo") != "he" {
		t.Errorf("expected 'he' but got '%s'", Delete("hello", "lo"))
	}
	if Delete("hello", "aeiou", "^e") != "hell" {
		t.Errorf("expected 'hell' but got '%s'", Delete("hello", "aeiou", "^e"))
	}
	if Delete("hello", "ej-m") != "ho" {
		t.Errorf("expected 'hell' but got '%s'", Delete("hello", "ej-m"))
	}
}

func Test_Insert(t *testing.T) {
	if result := Insert("foo", 1, "bar"); result != "fbaroo" {
		t.Errorf("expected '%s' but got '%s'", "fbaroo", result)
	}
	if result := Insert("foo", -2, "bar"); result != "fobaro" {
		t.Errorf("expected '%s' but got '%s'", "fobaro", result)
	}
	if result := Insert("foo", 0, "bar"); result != "barfoo" {
		t.Errorf("expected '%s' but got '%s'", "barfoo", result)
	}
}

func Test_InsertRune(t *testing.T) {
	str := "helloworld"
	out := InsertRune(str, ' ', 5)

	if out != "hello world" {
		t.Errorf("expected '%s' but got '%s'", "hello world", out)
	}

	out = InsertRune(str, '!', 10)

	if out != "helloworld!" {
		t.Errorf("expected '%s' but got '%s'", "helloworld!", out)
	}
}

func Test_InsertRunes(t *testing.T) {
	str := "helloworld"
	out := InsertRunes(str, 5, ',', ' ')

	if out != "hello, world" {
		t.Errorf("expected '%s' but got '%s'", "hello, world", out)
	}
}

func Test_Scrub(t *testing.T) {
	// Test for no repl given
	if result := Scrub("ab\uFFFDcd\xFF\xCEefg\xFF\xFC\xFD\xFAhijk"); result != "ab\uFFFDcd\uFFFDefg\uFFFDhijk" {
		t.Errorf("expected '%s' but got '%s'", "ab\uFFFDcd\uFFFDefg\uFFFDhijk", result)
	}
	// Test for repl function
	if result := Scrub("ab\uFFFDcd\xFF\xCEefg\xFF\xFC\xFD\xFAhijk", func(r rune) string {
		return "<i>"
	}); result != `ab<i>cd<i>efg<i>hijk` {
		t.Errorf("expected '%s' but got '%s'", `ab<i>cd<i>efg<i>hijk`, result)
	}
	// Test for empty (but present) repl
	if result := Scrub("ab\uFFFDcd\xFF\xCEefg\xFF\xFC\xFD\xFAhijk", ""); result != "abcdefghijk" {
		t.Errorf("expected '%s' but got '%s'", "abcdefghijk", result)
	}
}

func Test_Split(t *testing.T) {
	expected := []string{"now's", "the", "time"}

	expectMatchingSlices(t, "nil delimiter", " now's  the time ", nil, expected)

	expected = []string{"now's", "the time"}

	expectMatchingSlices(t, "nil delimiter with limit", " now's  the time ", nil, expected, 2)

	expected = []string{"now's", "the", "time"}

	expectMatchingSlices(t, "spacing prepended and whitespace delimiter", " now's  the time", " ", expected)

	expected = []string{"h", "e", "l", "l", "o"}

	expectMatchingSlices(t, "empty delimiter", "hello", "", expected)

	expected = []string{"Slip", "knot"}

	expectMatchingSlices(t, "static text for string and delimiter", "Slip:knot", ":", expected)

	expected = []string{"", "now's", "", "the", "time"}

	expectMatchingSlices(t, "spacing prepended and empty regexp", " now's  the time", regexp.MustCompile(` `), expected)

	expected = []string{"1", "2.34", "56", "7"}

	expectMatchingSlices(t, "regexp with whitespace following delimiter", "1, 2.34,56, 7", regexp.MustCompile(`,\s*`), expected)

	expected = []string{"h", "e", "l", "l", "o"}

	expectMatchingSlices(t, "empty regexp", "hello", regexp.MustCompile(``), expected)

	expected = []string{"h", "e", "llo"}

	expectMatchingSlices(t, "regexp with no content and a limit", "hello", regexp.MustCompile(``), expected, 3)

	expected = []string{"h", "i", "m", "o", "m"}

	expectMatchingSlices(t, "individual character split", "hi mom", regexp.MustCompile(`\s*`), expected)

	expected = []string{"m", "w y", "w"}

	expectMatchingSlices(t, "split in the middle of words", "mellow yellow", "ello", expected)

	expected = []string{"1", "2", "", "3", "4"}

	expectMatchingSlices(t, "split with no text between delimiters", "1,2,,3,4,,", ",", expected)

	expected = []string{"1", "2", "", "3,4,,"}

	expectMatchingSlices(t, "split with positive limit", "1,2,,3,4,,", ",", expected, 4)

	expected = []string{"1", "2", "", "3", "4", "", ""}

	expectMatchingSlices(t, "split with negative limit", "1,2,,3,4,,", ",", expected, -4)

	expected = []string{"1", ":", "", "", "2:3"}

	expectMatchingSlices(t, "regexp with groups", "1:2:3", regexp.MustCompile(`(:)()()`), expected, 2)

	expected = []string{}

	expectMatchingSlices(t, "empty source", "", ",", expected, -1)

	if errors == 1 {
		fmt.Printf("encountered 1 error during tests.\n")
	} else if errors > 0 {
		fmt.Printf("encountered %d errors during tests.\n", errors)
	}
}

func Test_Squeeze(t *testing.T) {
	if Squeeze("yellow moon", "") != "yelow mon" {
		t.Errorf("expected Squeeze(\"yellow moon\", \"\") to return \"yelow mon\" but got \"%s\"", Squeeze("yellow moon", ""))
	}
	if Squeeze("yellow moon") != "yelow mon" {
		t.Errorf("expected Squeeze(\"yellow moon\", \"\") to return \"yelow mon\" but got \"%s\"", Squeeze("yellow moon"))
	}
	if Squeeze("  now   is  the", " ") != " now is the" {
		t.Errorf("expected Squeeze(\"  now   is  the\", \"\") to return \" now is the\" but got \"%s\"", Squeeze("  now   is  the"))
	}
	if Squeeze("putters shoot balls", "m-z") != "puters shot balls" {
		t.Errorf("expected Squeeze(\"putters shoot balls\", \"m-z\") to return \"puters shot balls\" but got \"%s\"", Split("putters shoot balls", "m-z"))
	}
}

// expectMatchingSlices is a helper function that checks expected against actual slice returns
// and prints detailed errors if any. A hastily done helper for testing Split()
func expectMatchingSlices(t *testing.T, desc, str string, delim interface{}, expected []string, limit ...int) {
	var (
		errMsg string
		lim    int
		sts    string
	)

	result := Split(str, delim, limit...)

	if !equalSlices(result, expected) {
		st := getSplitType(delim, splitTypeString)
		if st == splitTypeAwk {
			sts = "AWK"
		} else if st == splitTypeString {
			sts = "String"
		} else if st == splitTypeChars {
			sts = "Chars"
		} else if st == splitTypeRegexp {
			sts = "Regexp"
		} else {
			sts = "Unknown"
		}

		fmt.Println("    Test:", strings.Title(desc))

		if limit != nil {
			lim = limit[0]
			errMsg = fmt.Sprintf("\nInput: \"%s\"\nDelimiter(%s): '%v'\nLimit: %d\nincorrect output:\n    Expected: %s (len %d)\n    Returned: %v (len %d)", str, sts, delim, lim, quoteSliceElements(expected), len(expected), quoteSliceElements(result), len(result))
		} else {
			errMsg = fmt.Sprintf("\nInput: \"%s\"\nDelimiter(%s): '%v'\nincorrect output:\n    Expected: %s (len %d)\n    Returned: %v (len %d)", str, sts, delim, quoteSliceElements(expected), len(expected), quoteSliceElements(result), len(result))
		}

		t.Error(errMsg)
		errors++
	}
}
