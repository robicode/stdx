# stdx/stringx: Extended String Library

This library contains extended functionality for strings (and slices thereof) that could prove useful, including many that are common in other languages.

Like the standard library, we assume all strings are UTF-8 encoded.

This package uses [Huan Du's xstrings library](https://github.com/huandu/xstrings) for much of the functionality, but this module tries to be compatible with the [Ruby language](https://ruby-lang.org) where possible, and thus many functions have been extended, (e.g. with default parameters or modes).

## License

MIT

## TODO & Contributing

Some things that remain:

* Get more functions converted e.g. `String#tr`, `String#count`, etc.
* Clean up the test suite
* More tests are always good
* Optimizations (maybe?)

I make no claims of being an expert at optimization, and there may well be evidence herein. If you find something that can be done better, a case where a function fails where it shouldn't, or any other issue, feel free to submit an issue (a PR is even better!) once I open up the issues tracker.
