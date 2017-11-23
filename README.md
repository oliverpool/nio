nio
==========

[![GoDoc](https://godoc.org/github.com/oliverpool/nio?status.svg)](https://godoc.org/github.com/oliverpool/nio)
[![Release](https://img.shields.io/github/release/oliverpool/nio.svg)](https://github.com/oliverpool/nio/releases/latest)
[![Build Status](https://travis-ci.org/oliverpool/nio.svg)](https://travis-ci.org/oliverpool/nio)
[![Go Report Card](https://goreportcard.com/badge/github.com/oliverpool/nio)](https://goreportcard.com/report/github.com/oliverpool/nio)

Usage
-----

The Buffer interface:

```go
type Buffer interface {
	io.ReadWriter
	Len() int
}
```

nio's Pipe method is a buffered version of io.Pipe
The writer return once its data has been written to the Buffer.
The reader returns with data off the Buffer.

```go
import "github.com/oliverpool/nio"

var buf bytes.Buffer
r, w := nio.NewPipe(&buf)
```


Licences
--------

The code in pipe.go is adapted from https://github.com/bradfitz/http2/blob/master/pipe.go
```
Copyright 2014 The Go Authors.
See https://code.google.com/p/go/source/browse/CONTRIBUTORS
Licensed under the same terms as Go itself:
https://code.google.com/p/go/source/browse/LICENSE
```

Some test and comments are adapted from https://github.com/oliverpool/nio:
```
The MIT License (MIT)

Copyright (c) 2015 Dustin H

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
the Software, and to permit persons to whom the Software is furnished to do so,
subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
```
