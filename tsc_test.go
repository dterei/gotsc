// Copyright 2016 David Terei.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// +build amd64

package gotsc

import (
	"testing"
	"time"
)

const (
	TSC_LOW  = 10
	TSC_HIGH = 200
)

func TestTSCOverhead(t *testing.T) {
	tsc := TSCOverhead()
	if tsc < 10 || tsc > 100 {
		t.Errorf("TSC Overhead returned number outside expected range: %d\n", tsc)
	}
}

func TestBench(t *testing.T) {
	tsc := TSCOverhead()
	start := BenchStart()
	end := BenchEnd()

	delta := end - start - tsc
	if start > end {
		t.Error("BenchEnd() earlier than BenchStart()")
	} else if delta > TSC_HIGH {
		t.Errorf("BenchEnd() - BenchStart() far greater than TSC overhead: %d\n",
			delta)
	}
}

func BenchmarkTime(b *testing.B) {
	for i := 0; i < b.N; i++ {
		time.Now()
	}
}

func BenchmarkBenchStart(b *testing.B) {
	for i := 0; i < b.N; i++ {
		BenchStart()
	}
}

func BenchmarkBenchEnd(b *testing.B) {
	for i := 0; i < b.N; i++ {
		BenchEnd()
	}
}
