# stdx/net/httpx: Extended net/http Library

This library consists of a small set of functions that may prove useful for development of HTTP applications. Currently there is not much here, but I plan to improve the library over time.

## Functions

### GetByteRanges

`GetByteRanges()` takes an HTTP `Range:` header and returns a valid set of `Range` objects. Invalid byte ranges are filtered out.

```go
type Range struct {
  From int64
  To   int64
}
```
### QValues

`QValues()` takes a header containing value-quality information and returns a set of value-quality pairs which can be sorted to obtain the best possible value. Returns a `[]QValue`.

```go
type QValue struct {
  Value   string
	Quality float64
}

// The QualityValues type is sortable via sort.Sort():
type QualityValues []QValue
```
### HTTPDate

`HTTPDate()` returns the time `t` formatted for use in HTTP headers. It really just calls `t.Format(http.TimeFormat)`.

## License

MIT

## Todo & Contributing

Some things that remain to be done:

* Clean up tests
* More tests
* More functions

I make no claims of being an expert at optimization, and there may well be evidence herein. If you find something that can be done better, a case where a function fails where it shouldn't, or any other issue, feel free to submit an issue (a PR is even better!) once I open up the issues tracker.