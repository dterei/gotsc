# gotsc

[![Build Status](https://travis-ci.org/dterei/gotsc.svg)](https://travis-ci.org/dterei/gotsc)
[![Go Report Card](https://goreportcard.com/badge/github.com/dterei/gotsc)](https://goreportcard.com/report/github.com/dterei/gotsc)
[![godoc](https://godoc.org/github.com/dterei/gotsc?status.svg)](http://godoc.org/github.com/dterei/gotsc)
[![BSD3 License](http://img.shields.io/badge/license-BSD3-brightgreen.svg?style=flat)][tl;dr Legal: BSD3]

[tl;dr Legal: BSD3]:
  https://tldrlegal.com/license/bsd-3-clause-license-(revised)
  "BSD3 License"

Golang library for access the CPU timestamp cycle counter (TSC) on x86-64. If
not familar with using the `TSC` for benchmarking, refer to the [Intel
whitepaper](http://www.intel.com/content/www/us/en/embedded/training/ia-32-ia-64-benchmark-code-execution-paper.html).

Golang 1.4 or later is currently supported and x86-64 architetcture. The
package will build on other architectures but all functions will simply return
0.

## Compared with time.Now()

There are two advantages over the standard golang `time.Now()` function:

1) Measurement is in cycles - for many situations cycle count is a more
   informative number than wall-clock time.
2) Lower overhead - `time.Now()` takes ~32K cycles to call while `BenchStart()`
   and `BenchEnd()` take ~250 cycles.

## Usage

``` .go
package main

import (
  "fmt"
  "github.com/dterei/gotsc"
)

const N = 100

func main() {
  tsc := gotsc.TSCOverhead()
  fmt.Println("TSC Overhead:", tsc)

  var avg uint64
  for i := 0; i < N; i++ {
    start := gotsc.BenchStart()
    // code to evaluate
    end := gotsc.BenchEnd()
    avg += end - start - tsc
  }

  fmt.Println("Cycles:", avg / N)
}
```

## Licensing

This library is BSD-licensed.

## Get involved!

We are happy to receive bug reports, fixes, documentation enhancements,
and other improvements.

Please report bugs via the
[github issue tracker](http://github.com/dterei/gotsc/issues).

Master [git repository](http://github.com/dterei/gotsc):

* `git clone git://github.com/dterei/gotsc.git`

## Authors

This library is written and maintained by David Terei, <code@davidterei.com>.

