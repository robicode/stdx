# stdx/stringx Wanted Functionality

This file contains functions I would like to have added to the library, but either don't have time to implement, are too difficult for me to implement, or I just don't see how they can be converted to Go's syntax and semantics.

If/when they are implemented, I will move their section below into the README and remove them here.

## Functions

Note that I am aware some of these functions are available in either the standard library or `xtrings`, but those versions lack features documented herein.

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
Chop("hello \n there")       // "hello \n there"
Chomp("hello", "llo")        // "he"
Chomp("hello\r\n\r\n", "")   // "hello"
Chomp("hello\r\n\r\r\n", "") // "hello\r\n\r"
```

### Unpack

Decodes `str` (which may contain binary data) according to the
`format` string, returning a slice of each value extracted. The `format`
string consists of a sequence of single-character directives, summarized
in the table at the end of this entry. Each directive may be followed by
a number, indicating the number of times to repeat with this directive.
An asterisk (*) will use up all remaining elements. The directives
sSiIlL may each be followed by an underscore (_) or exclamation mark
(!) to use the underlying platform's native size for the specified
type; otherwise, it uses a platform-independent consistent size. Spaces
are ignored in the `format` string.

**Not Implemented Because:** Not sure if I can port this yet due to difficulty, lack of time, and/or lack of Go features. If I do port it, some field types in the format table probably won't be supported.

```go
Unpack("abc \0\0abc \0\0", "A6Z6")          // []string{"abc", "abc "}
Unpack("abc \0\0", "a3a3")                  // []string{"abc", " \000\000"}
Unpack("abc \0abc \0", "Z*Z*")              // []string{"abc ", "abc "}
Unpack("aa", "b8B8")                        // []string{"10000110", "01100001"}
Unpack("aaa", "h2H2c")                      // []string{"16", "61", "97"}
Unpack("\xfe\xff\xfe\xff", "sS")            // []string{"-2", "65534"}
Unpack("now=20is", "M*")                    // []string{"now is"}
Unpack("whole", "xax2aX2aX1aX2a")           // []string{"h", "e", "l", "l", "o"}
```

The table below summarizes the various formats. Note that unlike the Ruby version of this function, we always return a slice of []string, because as far as I know Go does not support multiple variable types in the same slice.

```
  Integer       |         |
  Directive     | Returns | Meaning
  ------------------------------------------------------------------
  C             | Integer | 8-bit unsigned (unsigned char)
  S             | Integer | 16-bit unsigned, native endian (uint16_t)
  L             | Integer | 32-bit unsigned, native endian (uint32_t)
  Q             | Integer | 64-bit unsigned, native endian (uint64_t)
  J             | Integer | pointer width unsigned, native endian (uintptr_t)
                |         |
  c             | Integer | 8-bit signed (signed char)
  s             | Integer | 16-bit signed, native endian (int16_t)
  l             | Integer | 32-bit signed, native endian (int32_t)
  q             | Integer | 64-bit signed, native endian (int64_t)
  j             | Integer | pointer width signed, native endian (intptr_t)
                |         |
  S_ S!         | Integer | unsigned short, native endian
  I I_ I!       | Integer | unsigned int, native endian
  L_ L!         | Integer | unsigned long, native endian
  Q_ Q!         | Integer | unsigned long long, native endian (ArgumentError
                |         | if the platform has no long long type.)
  J!            | Integer | uintptr_t, native endian (same with J)
                |         |
  s_ s!         | Integer | signed short, native endian
  i i_ i!       | Integer | signed int, native endian
  l_ l!         | Integer | signed long, native endian
  q_ q!         | Integer | signed long long, native endian (ArgumentError
                |         | if the platform has no long long type.)
  j!            | Integer | intptr_t, native endian (same with j)
                  |         |
  S> s> S!> s!> | Integer | same as the directives without ">" except
  L> l> L!> l!> |         | big endian
  I!> i!>       |         |
  Q> q> Q!> q!> |         | "S>" is same as "n"
  J> j> J!> j!> |         | "L>" is same as "N"
                |         |
  S< s< S!< s!< | Integer | same as the directives without "<" except
  L< l< L!< l!< |         | little endian
  I!< i!<       |         |
  Q< q< Q!< q!< |         | "S<" is same as "v"
  J< j< J!< j!< |         | "L<" is same as "V"
                |         |
  n             | Integer | 16-bit unsigned, network (big-endian) byte order
  N             | Integer | 32-bit unsigned, network (big-endian) byte order
  v             | Integer | 16-bit unsigned, VAX (little-endian) byte order
  V             | Integer | 32-bit unsigned, VAX (little-endian) byte order
                |         |
  U             | Integer | UTF-8 character
  w             | Integer | BER-compressed integer (see Array#pack)

  Float        |         |
  Directive    | Returns | Meaning
  -----------------------------------------------------------------
  D d          | Float   | double-precision, native format
  F f          | Float   | single-precision, native format
  E            | Float   | double-precision, little-endian byte order
  e            | Float   | single-precision, little-endian byte order
  G            | Float   | double-precision, network (big-endian) byte order
  g            | Float   | single-precision, network (big-endian) byte order

  String       |         |
  Directive    | Returns | Meaning
  -----------------------------------------------------------------
  A            | String  | arbitrary binary string (remove trailing nulls and ASCII spaces)
  a            | String  | arbitrary binary string
  Z            | String  | null-terminated string
  B            | String  | bit string (MSB first)
  b            | String  | bit string (LSB first)
  H            | String  | hex string (high nibble first)
  h            | String  | hex string (low nibble first)
  u            | String  | UU-encoded string
  M            | String  | quoted-printable, MIME encoding (see RFC2045)
  m            | String  | base64 encoded string (RFC 2045) (default)
               |         | base64 encoded string (RFC 4648) if followed by 0
  P            | String  | pointer to a structure (fixed-length string)
  p            | String  | pointer to a null-terminated string

  Misc.        |         |
  Directive    | Returns | Meaning
  -----------------------------------------------------------------
  @            | ---     | skip to the offset given by the length argument
  X            | ---     | skip backward one byte
  x            | ---     | skip forward one byte
```
