# streampos: a writer that tracks lines and columns

Or *"How I learned to love composition even more!"*

[![GoDoc](https://godoc.org/github.com/phf/streampos?status.svg)](https://godoc.org/github.com/phf/streampos)
[![Go Report Card](https://goreportcard.com/badge/github.com/phf/streampos)](https://goreportcard.com/report/github.com/phf/streampos)

## Summary

Package `streampos` exports an `io.Writer` that tracks (line, column)
positions in streams, suitable for error messages. It looks for newline
characters and builds a mapping from byte offset ranges (which start at
0) to line numbers (which start at 1). The `Position` method returns the
line and column number that corresponds to a given offset in the stream.
The offset has to be between 0 and whatever the `Length` method returns.

If all of this strikes you as rather strange, you should probably read
my blog post once I am done with it. (Watch this space for the link.)

## Usage

You should be able to just say this:

	go get -u github.com/phf/streampos

Then you can import and use the package as follows:

```golang
package main

import (
	"fmt"
	"github.com/phf/streampos"
)

func main() {
	b := []byte("Write\nmore\nGo!\n")

	w := &streampos.Writer{}
	w.Write(b) // never fails

	for i := int64(0); i < int64(len(b)); i++ {
		l, c, err := w.Position(i)
		fmt.Printf("offset %v maps to line %v, column %v (error: %v)\n", i, l, c, err)
	}
}
```

Enjoy!

## TODO

Column numbers are currently computed in terms of *bytes* just like in the
Go compiler. The advantage is a clear definition of what "column" actually
means. The drawback is that users may not agree when it comes to TAB
characters or multi-byte runes. However, given that GCC's bug-tracker is
*full* of column-related issues, maybe it's better to keep it simple?

Maybe we could at least make column numbers *somewhat* configurable. The
`text/scanner` package for example uses character counts instead of byte
counts, at least that's how I understand the documentation. Adding a small
constructor function with an (optional) parameter to specify options could
achieve that in a moderately straightforward way. I don't care very much
one way or the other, I am happy with columns in terms of bytes.
At least for now.

Currently I use a *slice* of byte offset ranges and a *linear* search.
For large streams, that's obviously slow. Binary search would help, as
would some tree-based data-structure that supports range queries. So far
I don't care enough for either, but I'd be happy to accept a pull-request.

Concurrency safety, anyone? :-)

## License

The MIT License.
