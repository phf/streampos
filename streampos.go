// Copyright 2016 Peter H. Froehlich. All rights reserved.
// Use of this source code is governed by the MIT license,
// see the LICENSE.md file.

// Package streampos exports a Writer that tracks positions (line, column)
// suitable for error messages.
package streampos

import (
	"fmt"
)

// position maps a byte offset range to a line number.
type position struct {
	from, to int64
	line     int64
}

// Writer looks for newline characters and constructs a mapping from byte
// offset ranges to line numbers. It is typically used with a TeeReader to
// track positions in some input stream.
type Writer struct {
	positions []position // byte range <-> line number
	offset    int64      // byte offset after last Write
	line      int64      // line number after last Write
	total     int64      // total bytes after last Write
}

// Write accepts data to track positions in. It does not perform any actual
// I/O, so writes are never short.
func (w *Writer) Write(b []byte) (n int, err error) {
	s := string(b)
	// offset and line for *this* Write call
	var o int64
	var l int64

	for i, r := range s {
		if r == '\n' {
			pos := position{w.offset + o, w.offset + int64(i), w.line + l + 1}
			w.positions = append(w.positions, pos)
			o = int64(i) + 1
			l++
		}
	}

	// update offset and line for the *next* Write call
	w.offset += o
	w.line += l
	// update total to allow for queries after last byte range
	w.total += int64(len(b))

	return len(b), nil
}

// Length returns how many bytes have been processed so far.
func (w *Writer) Length() int64 {
	return w.total
}

// Line returns the line number for the given offset.
// See Position for more information.
func (w *Writer) Line(offset int64) (int64, error) {
	line, _, err := w.Position(offset)
	return line, err
}

// Column returns the column number for the given offset.
// See Position for more information.
func (w *Writer) Column(offset int64) (int64, error) {
	_, column, err := w.Position(offset)
	return column, err
}

// Position returns the line and column number for the given offset in the
// underlying data stream. The offset cannot be negative but must be less
// than Length. Line numbers start at 1 from the beginning of the stream,
// column numbers start at 1 from the beginning of the line (left to right).
// Note that column numbers are based on bytes, so '\t' counts as 1 column
// whereas multi-byte runes count as several.
func (w *Writer) Position(offset int64) (line, column int64, err error) {
	// validate offset
	if 0 > offset || offset >= w.total {
		return -1, -1, fmt.Errorf("streampos: offset %v out of range, must be in [0, %v]", offset, w.total-1)
	}
	// look for the proper range
	for _, p := range w.positions {
		if p.from <= offset && offset <= p.to {
			return p.line, offset - p.from + 1, nil
		}
	}
	// we could have seen more bytes but not formed a range since
	// there was no \n yet
	p := w.positions[len(w.positions)-1]
	if p.to < offset && offset < w.total {
		return p.line + 1, offset - p.to, nil
	}
	return -1, -1, fmt.Errorf("streampos: internal error")
}
