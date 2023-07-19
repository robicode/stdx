# stdx/stringx: Extended String Library

This library contains extended functionality for strings (and slices thereof) that could prove useful, including many that are common in other languages.

Like the standard library, we assume all strings are UTF-8 encoded.

This package uses [Huan Du's xstrings library](https://github.com/huandu/xstrings) for much of the functionality, but this module tries to be compatible with the [Ruby language](https://ruby-lang.org) where possible, and thus many functions have been extended, (e.g. with default parameters or modes of operation).

Some functionality is also pulled from [Mawuli Kofi Adzoe's stringutil](https://github.com/wallclockbuilder/stringutil) package and modified for this library.

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

### Chomp

Returns a new String with the given record separator removed from the
end of str (if present). If using the default separator, then chomp also 
removes carriage return characters (that is it will remove \n, \r, and \r\n).
If $/ is an empty string, it will remove all trailing newlines from the string.

```go
Chomp("hello")               // "hello"
Chomp("hello\n")             // "hello"
Chomp("hello\r\n")           // "hello"
Chomp("hello\n\r")           // "hello\n"
Chomp("hello\r")             // "hello"
Chomp("hello \n there")      // "hello \n there"
Chomp("hello", "llo")        // "he"
Chomp("hello\r\n\r\n", "")   // "hello"
Chomp("hello\r\n\r\r\n", "") // "hello\r\n\r"
```

### Chr

Chr returns a string containing the first rune of str.

```go
Chr("foo")    // "f"
```

### Count

Each `other_str` parameter defines a set of characters to count.  The
intersection of these sets defines the characters to count in `str`.  Any
`other_str` that starts with a caret (^) is negated.  The sequence `c1-c2`
means all characters between `c1` and `c2`.  The backslash character ("\") can
be used to escape ^ or - and is otherwise ignored unless it appears at
the end of a sequence or the end of a `other_str`.

**Not Implemented Because:** The version in `xtrings` has a different algorithm for its pattern. It is just as valid a function, but I'm trying to stick with Ruby algorithms here. The version in stdlib doesn't support patterns at all.

```go
str := "hello world"
Count(str, "lo")                 // 5
Count(str, "lo", "o")            // 2
Count(str, "hello", "^l")        // 4
Count("hello", "ej-m")           // 4

str = "hello^world"
Count(str, "\\^aeiou")           // 4
Count(str, "a\\-eo")             // 3

str = "hello world\\r\\n"
Count(str, "\\")                 // 2
Count(str, "\\A")                // 0
Count(str, "A-\\w")              // 3
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

### DeleteMatchingRunes

DeleteMatchingRunes deletes all instances of the rune r from the given
string str.

```go
DeleteMatchingRunes("hello", 'l')        // "heo"
```

### DeleteRune

DeleteRune deletes a single rune from string str
at the given index.

```go
DeleteRune("hello", 3)        // "helo"
```

Note that this may not always correspond to the actual byte
index as UTF-8 runes may span multiple bytes.

### FormatStrings

`FormatStrings` takes a []string and returns a string similar to that
returned by `fmt.Println([]string{...})`, but with each element quoted for
readability.

```go
slice := []string{"one", "two", "three", "four "}

// Normal Output
fmt.PrintLn(slice)                   // [one two three four ]

// With FormatStrings:
fmt.Println(FormatStrings(slice))    // ["one", "two", "three", "four "]
```

### Gsub

`Gsub` Returns a copy of str with all occurrences of
pattern substituted for the second argument. The
pattern is typically a Regexp; if given as a String, any
regular expression metacharacters it contains will be interpreted
literally, e.g. \d will match a backslash followed by 'd', instead of a
digit.

If replacement is a string it will be substituted for the matched text.
It may contain back-references to the pattern's capture groups of the
form `\d`, where d is a group number. Unlike in Ruby, `\k` is unsupported.

```go
GSub("hello", regexp.MustCompile(`[aeiou]`), "*")                 // "h*ll*"
GSub("hello", regexp.MustCompile(`([aeiou])`), "<\1>")            // "h<e>ll<o>"
GSub("hello", regexp.MustCompile(`.`), func(s string) string {    // "104 101 108 108 111 "
	return strconv.Itoa(Ord(s)) + " "
})
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

### Scan

Both forms iterate through str, matching the pattern (which may be
a `*Regexp` or a string). For each match, a result is generated and either
added to the result array or passed to the function. If the pattern
contains no groups, each individual result consists of the matched
string.  If the pattern contains groups, each individual result is
itself a slice containing one entry per group.

```go
a := "cruel world"
Scan(a, regexp.MustCompile(`\w+`))        // ["cruel", "world"]
Scan(a, regexp.MustCompile(`...`))        // ["cru", "el ", "wor"]
Scan(a, regexp.MustCompile(`(...)`))      // [["cru"], ["el "], ["wor"]]
Scan(a, regexp.MustCompile(`(..)(..)`))   // [["cr", "ue"], ["l ", "wo"]]
```
And when given a function:

```go
Scan(a, regexp.MustCompile(`\w+`), func(match interface{}){
	...
	fmt.Printf("<<%s>> \n", match)
})
```

produces:

```
<<cruel>> <<world>>
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

### Tr

Returns a copy of str with the characters in from_str replaced by the
corresponding characters in to_str.  If to_str is shorter than from_str,
it is padded with its last character in order to maintain the
correspondence.

```go
Tr("hello", "el", "ip")        // "hippo"
Tr("hello", "aeiou", "*")      // "h*ll*"
Tr("hello", "aeiou", "AA*")    // "hAll*"
```

Both strings may use the c1-c2 notation to denote ranges of characters,
and from_str may start with a ^, which denotes all characters except
those listed.

```go
Tr("hello", "a-y", "b-z")     // "ifmmp"
Tr("hello", "^aeiou", "*")    // "*e**o"
```

The backslash character \ can be used to escape ^ or - and is otherwise
ignored unless it appears at the end of a range or the end of the
from_str or to_str:

```go
Tr("hello^world", "\\^aeiou", "*") // "h*ll**w*rld"
Tr("hello-world", "a\\-eo", "*")   // "h*ll**w*rld"

Tr("hello\r\nworld", "\r", "")     // "hello\nworld"
Tr("hello\r\nworld", "\\r", "")    // "hello\r\nwold"
Tr("hello\r\nworld", "\\\r", "")   // "hello\nworld"

Tr("X['\\b']", "X\\", "")          // "['b']"
Tr("X['\\b']", "X-\\]", "")        // "'b'"
```

## License

This package is licensed under MIT.

However, the code does use two external libraries:

* [Huan Du's xstrings library](https://github.com/huandu/xstrings) is used for some functionality, and it is also licensed under MIT. Please see the LICENSE file for details.
* [Mawuli Kofi Adzoe's stringutil](https://github.com/wallclockbuilder/stringutil) package is used for some functionality, and it is licensed under a BSD 2-clause license. See the LICENSE file for details.

## TODO & Contributing

Here are some things that could be done to improve the library:

* Get more functions, particularly those in the [wish list](WISHLIST.md) ported.
* Additional useful functions that aren't in other languages?
* Clean up the test suite
* More tests are always good
* Optimizations (maybe?)

I make no claims of being an expert at optimization, and there may well be evidence herein. If you find something that can be done better, a case where a function fails where it shouldn't, or any other issue, feel free to submit an issue (a PR is even better!) once I open up the issues tracker.
