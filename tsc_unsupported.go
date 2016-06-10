// Copyright 2016 David Terei.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// +build !amd64

// Package gotsc provides access to the timestamp cycle counter on x86-64 for
// performing close to cycle accurate benchmarking on x86-64. All recent (think
// since 2010) generation Intel CPU's provide a global, synchronized cycle
// counter great for benchmarking and time measurement across all cores.
//
// On non x86-64 platforms, all functions simply return 0 at this time.
//
package gotsc

// BenchStart obtains the cycle counter. It should be used at the start of
// benchmarking some code.
func BenchStart() uint64 {
	return 0
}

// BenchEnd obtains the cycle counter. It should be used at the end of
// benchmarking some code. There is a subtle difference in the allowed
// reordering of operations between BenchEnd and BenchStart, hence the two
// functions rather than one.
func BenchEnd() uint64 {
	return 0
}

// TSCOverhead measures the cycle overhead of calling the underlying `rdtsc`
// instruction to obtain the current CPU cycle count. You should subtract this
// value from all your cycle count measurements to accurately benchmark some
// code.
func TSCOverhead() uint64 {
	return 0
}
