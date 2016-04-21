// Copyright 2016 Peter H. Froehlich. All rights reserved.
// Use of this source code is governed by the MIT license,
// see the LICENSE.md file.

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
