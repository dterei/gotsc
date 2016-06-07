// Package gotsc provides access to the timestamp cycle counter on x86-64 for
// performing close to cycle accurate benchmarking on x86-64. All recent (think
// since 2010) generation Intel CPU's provide a global, synchronized cycle
// counter great for benchmarking and time measurement across all cores.
package gotsc

// BenchStart obtains the cycle counter. It should be used at the start of
// benchmarking some code.
func BenchStart() uint64

// BenchEnd obtains the cycle counter. It should be used at the end of
// benchmarking some code. There is a subtle difference in the allowed
// reordering of operations between BenchEnd and BenchStart, hence the two
// functions rather than one.
func BenchEnd() uint64

// TSCOverhead measures the cycle overhead of calling the underlying `rdtsc`
// instruction to obtain the current CPU cycle count. You should subtract this
// value from all your cycle count measurements to accurately benchmark some
// code.
func TSCOverhead() uint64 {
	var t0, t1 uint64
	overhead := uint64(1000000000000000000)

	for i := 0; i < 100000; i++ {
		t0 = BenchStart()
		t1 = BenchEnd()
		if t1-t0 < overhead {
			overhead = t1 - t0
		}
	}

	return overhead
}
