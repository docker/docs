// Copyright 2013 The Go-SQLite Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package codec

import (
	"encoding/hex"
	"fmt"
	"os"
	"strconv"

	. "code.google.com/p/go-sqlite/go1/sqlite3"
)

type hexDump struct {
	key   []byte
	out   *os.File
	quiet bool
	res   int
}

func newHexDump(ctx *CodecCtx, key []byte) (Codec, *Error) {
	_, opts, tail := parseKey(key)
	c := &hexDump{key, os.Stderr, false, -1}

	// Set options
	for k, v := range opts {
		switch k {
		case "quiet":
			c.quiet = true
		case "reserve":
			if n, err := strconv.ParseUint(v, 10, 8); err == nil {
				c.res = int(n)
			}
		default:
			return nil, NewError(MISUSE, "invalid codec option: "+k)
		}
	}

	// Open output file
	switch file := string(tail); file {
	case "":
	case "-":
		c.out = os.Stdout
	default:
		var err error
		c.out, err = os.OpenFile(file, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
		if err != nil {
			return nil, NewError(ERROR, err.Error())
		}
	}

	fmt.Fprintf(c.out, "Init (\n"+
		"    Path=%s\n"+
		"    Name=%s\n"+
		"    PageSize=%d\n"+
		"    Reserve=%d\n"+
		"    Fixed=%t\n"+
		")\n",
		ctx.Path, ctx.Name, ctx.PageSize, ctx.Reserve, ctx.Fixed)
	return c, nil
}

func (c *hexDump) Reserve() int {
	fmt.Fprintf(c.out, "Reserve\n")
	return c.res
}

func (c *hexDump) Resize(pageSize, reserve int) {
	fmt.Fprintf(c.out, "Resize (pageSize=%d, reserve=%d)\n", pageSize, reserve)
}

func (c *hexDump) Encode(page []byte, pageNum uint32, op int) ([]byte, *Error) {
	fmt.Fprintf(c.out, "Encode (pageNum=%d, op=%d)\n", pageNum, op)
	c.dump(page)
	return page, nil
}

func (c *hexDump) Decode(page []byte, pageNum uint32, op int) *Error {
	fmt.Fprintf(c.out, "Decode (pageNum=%d, op=%d)\n", pageNum, op)
	c.dump(page)
	return nil
}

func (c *hexDump) Key() []byte {
	fmt.Fprintf(c.out, "Key\n")
	return c.key
}

func (c *hexDump) Free() {
	fmt.Fprintf(c.out, "Free\n")
	if c.out != os.Stdout && c.out != os.Stderr {
		c.out.Close()
	}
}

func (c *hexDump) dump(b []byte) {
	if !c.quiet {
		hd := hex.Dumper(c.out)
		hd.Write(b)
		hd.Close()
		c.out.Write([]byte("\n"))
	}
}
