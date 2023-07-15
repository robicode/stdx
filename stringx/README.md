# stdx/stringx: Extended String Library

This library contains extended functionality for strings (and slices thereof) that could prove useful, including many that are common in other languages.

Like the standard library, we assume all strings are UTF-8 encoded.

This package uses [Huan Du's xstrings library](https://github.com/huandu/xstrings) for much of the functionality, but this module tries to be compatible with the [Ruby language](https://ruby-lang.org) where possible, and thus many functions have been extended, (e.g. with default parameters or modes of operation).

Where a function can be found in other languages' standard libraries, I try to stay close to the Ruby syntax and features, as I am most familiar with that language and the features of its standard library are quite robust.

## Functions
### CaseCmp

Compares self and other string, ignoring case, and returns -1 if other string is larger, 0 if the two are equal, or -1 if other string is smaller.

```go
CaseCmp("foo", "foo")        // 0
CaseCmp("foo", "food")       // -1
CaseCmp("food", "foo")       // 1
CaseCmp("FOO", "foo")        // 0
CaseCmp("foo", "FOO")        // 0
```

### Center

Centers `str` in `width`.  If `width` is greater than the length of `str`,
returns a new string of length `width` with `str` centered and padded with
`pad`; otherwise, returns `str`.

```go
Center("hello", 4)            // "hello"
Center("hello", 20)           // "       hello        "
Center("hello", 20, "123")    // "1231231hello12312312"
```

### Delete

Returns a copy of `str` with all characters in the intersection of
its arguments deleted. Uses the same rules for building the set of
characters as `Count()`.

```go
Delete("hello", "l", "lo")            // "yelow mon"
Delete("hello", "lo")                 // "he"
Delete("hello", "aeiou", "^e")        // "hell"
Delete("hello", "ej-m")               // "ho"
```

### FormatStrings

FormatStrings takes a []string and returns a string similar to that
returned by `fmt.Println([]string{...})`, but with each element quoted for
readability.

```go
slice := []string{"one", "two", "three", "four "}

// Normal Output
fmt.PrintLn(slice)                   // [one two three four ]

// With FormatStrings:
fmt.Println(FormatStrings(slice))    // ["one", "two", "three", "four "]
```

### Insert

`Insert` Inserts the given `other` string into `str`; returns the new string.

If the Integer `index` is positive or zero, inserts `other` at offset `index`:

```go
Insert("foo", 1, "bar")            // "fbaroo"
```

If the Integer `index` is negative, counts backward from the end of `str`
and inserts `other` at offset `index + 1` (that is, after
`str[index])`:

```go
Insert("foo", -2, "bar")           // "fobaro"
```

Insert does no utf-8 validation of the `other` string.

### InsertRune

`InsertRune` inserts rune `r` into `s` at position `p`. No checking is done to validate `r` is a valid utf-8 `rune`.

If `p` is greater than `len(s) - 1`, it is set to `len(p) - 1`.

```go
str := "helloworld"
InsertRune(str, ' ', 5)                 // "hello world"
InsertRune(str, '!', 10)                // "helloworld!"
```

### InsertRunes

`InsertRunes` is like `InsertRune`, but you can pass multiple runes
and they will all be inserted at the given position `p`.

If `p` is greater than `len(s) - 1`, it is set to `len(p) - 1`.

```go
str := "helloworld"
InsertRunes(str, 5, ',', ' ')                 // "hello, world"
```

### IsASCII

IsASCII returns true if s consists entirely of ASCII characters. Pulled straight from stdlib and exported.

### Partition

Partition Searches sep or pattern (regexp) in the string and
returns the part before it, the match, and the part after it. If it is
not found, returns two empty strings and str.

```go
Partition("hello", "l")                      // []string{"he", "l", "lo"}
Partition("hello", "x")                      // []string{"hello", "", ""}
Partition("hello", regexp.MustCompile(`.l`)) // []string{"h", "el", "lo"}
```
### Scrub

If the string contains any invalid byte sequences then replace invalid
bytes with given replacement string, else returns `str`. If `repl` is given
as a function, replace invalid bytes with returned value of the
function.

```go
str := "ab\uFFFDcd\xFF\xCEefg\xFF\xFC\xFD\xFAhijk"
Scrub(str)                        // "ab\uFFFDcd\uFFFDefg\uFFFDhijk"
Scrub(str, "")                    // "abcdefghijk"
Scrub(str, ".")                   // "ab.cd.efg.hijk"

Scrub(str, func(r rune) string {  // "ab!cd!efg!hijk"
	return "!"
})
```

### Split
`Split` divides `s` into substrings based on a `pattern`, returning a
slice of these substrings.

If `pattern` is a `string`, then its contents are used as the
delimiter when splitting `s`. If `pattern` is a single
space, `s` is split on whitespace, with leading and trailing
whitespace and runs of contiguous whitespace characters ignored.

If `pattern` is a `Regexp`, `s` is divided where the `pattern` matches. Whenever
the `pattern` matches a zero-length `string`, `s` is split into individual characters. 
If `pattern` contains groups, the respective matches will be returned in the slice as well.

If `pattern` is nil, `s` is split on whitespace as if " " had been passed
as the `pattern`.

If the `limit` parameter is omitted, trailing null fields are
suppressed. If `limit` is a positive number, at most that number
of split substrings will be returned (captured groups will be returned
as well, but are not counted towards the `limit`). If `limit` is
1, the entire string is returned as the only entry in a slice. If
negative, there is no limit to the number of fields returned, and
trailing null fields are not suppressed.

When the input `s` is empty an empty slice is returned as the `string` is
considered to have no fields to split.

```go
Split(" now's  the time ")                 // []string{"now's", "the", "time"}
Split(" now's  the time ", ' ')            // []string{"now's", "the", "time"}

re := regexp.MustCompile(` `)
Split(" now's  the time", re)              // []string{"", "now's", "", "the", "time"}

re = regexp.MustCompile(`,\s*`)
Split("1, 2.34,56, 7", re)                 // []string{"1", "2.34", "56", "7"}

re = regexp.MustCompile(``)
Split("hello", re)                         // []string{"h", "e", "l", "l", "o"}
Split("hello", re, 3)                      // []string{"h", "e", "llo"}

re = regexp.MustCompile(`\s*`)
Split("hi mom", re))                       // []string{"h", "i", "m", "o", "m"}

Split("mellow yellow", "ello")             // []string{"m", "w y", "w"}
Split("1,2,,3,4,,", ",")                   // []string{"1", "2", "", "3", "4"}
Split("1,2,,3,4,,", ",", 4)                // []string{"1", "2", "", "3,4,,"}
Split("1,2,,3,4,,", ",", -4)               // []string{"1", "2", "", "3", "4", "", ""}

re = regexp.MustCompile(`(:)()()`)
Split("1:2:3", re, 2)                      // []string{"1", ":", "", "", "2:3"}

Split("", ",", -1)                         // []string{}
```

### Squeeze

`Squeeze` builds a set of characters from the `other` string
parameter(s) using the procedure described for `Count()`. Returns a
new string where runs of the same character that occur in this set are
replaced by a single character. If no arguments are given, all runs of
identical characters are replaced by a single character.

```go
Squeeze("yellow moon", "")             // "yelow mon"
Squeeze("  now   is  the", "m-z")      // " now is the"
Squeeze("putters shoot balls", " ")    // "puters shot balls"
```

## License

MIT

## TODO & Contributing

Here are some things that could be done to improve the library:

* Get more functions, particularly those in the [wish list](WISHLIST.md) ported.
* Additional useful functions that aren't in other languages?
* Clean up the test suite
* More tests are always good
* Optimizations (maybe?)

I make no claims of being an expert at optimization, and there may well be evidence herein. If you find something that can be done better, a case where a function fails where it shouldn't, or any other issue, feel free to submit an issue (a PR is even better!) once I open up the issues tracker.
