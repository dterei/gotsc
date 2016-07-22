# gotsc

[![Build Status](https://travis-ci.org/dterei/gotsc.svg)](https://travis-ci.org/dterei/gotsc)
[![Go Report Card](https://goreportcard.com/badge/github.com/dterei/gotsc)](https://goreportcard.com/report/github.com/dterei/gotsc)
[![godoc](https://godoc.org/github.com/dterei/gotsc?status.svg)](http://godoc.org/github.com/dterei/gotsc)
[![BSD3 License](http://img.shields.io/badge/license-BSD3-brightgreen.svg?style=flat)][tl;dr Legal: BSD3]

[tl;dr Legal: BSD3]:
  https://tldrlegal.com/license/bsd-3-clause-license-(revised)
  "BSD3 License"

Golang library for access the CPU timestamp cycle counter (TSC) on x86-64. If
not familar with using the TSC for benchmarking, refer to the
[Intel whitepaper][intel1]. This is designed to be used for benchmarking code, so
takes steps to prevent instruction reordering across measurement boundaries by
the CPU.

Golang 1.4 or later is currently supported and x86-64 architetcture. The
package will build on other architectures but all functions will simply return
0.

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

  start := gotsc.BenchStart()
  for i := 0; i < N; i++ {
    // code to evaluate
  }
  end := gotsc.BenchEnd()
  avg := (end - start - tsc) / N

  fmt.Println("Cycles:", avg)
}
```

## Compared with time.Now()

There are two advantages over the standard golang `time.Now()` function:

1. Measurement is in cycles - for many situations cycle count is a more
   informative number than wall-clock time.
2. Careful use of CPU serializing instructions to ensure no code you are
   benchmarking is moved outside the timed region, and no code you aren't
   benchmarking is moved into it.

Claim (2) may be a little contensious, so see below. For benchmarking with the
TSC, we use the approach suggested by [Intel][intel1]:

``` .asm
cpuid
rdtsc
// code to benchmark
rdtscp
cpuid
```

## Reading the TSC

There appears to be only confusion on what is both the correct and best way to
read the TSC when benchmarking code. The most obvious and naive approach would
be:

``` .asm
// code before
rdtsc
// code to benchmark
rdtsc
// code after
```

But `rdtsc` doesn't prevent instructions being reordered by the CPU around it.
Thus, code before and after could move into the benchmarked region, while code
in the benchmarked region could move out.

The best Intel documentation on this suggests the following approach:

``` .asm
// code before
cpuid
rdtsc
// code to benchmark
rdtscp
cpuid
/// code after
```

The `cpuid` instruction is a full barrier, preventing reordering in both
directions, while `rdtscp` prevents reordering from above. We use `rdtscp` at
the end rather than `cpuid; rdtsc` as `cpuid` is an expensive instruction with
high variance, so we want it outside the benchmarked region.

Ideally our benchmarking approach provides the following:

1. Low variance for instructions involved in retrieving start and end TSC so
   that we can subtract their overhead from the measurement with more
   confidence.
2. Low cost to read the TSC so that we can take benchmarks as often as possible
   without affecting application performance.
3. High resolution so that we can measure the cost of very small sets of
   instructions.

The recommended Intel approach provides (1) and (3) but the use of `cpuid` is
fairly expensive. The Intel SDM suggest that the `lfence` instruction can be
used as an alternative to `cpuid`, while AMD suggest the use of `mfence`.

Linux takes this approach. [Originally][lxr1] ([LKML][lkml1]), it used an
`lfence` either side with the thinking being that `lfence` only prevents
reordering from above:

``` .asm
// Linux kernel TSC usage (circa 2008)
lfence
rdtsc
lfence
```

This was [later][lxr2] ([LKML][lkml2]) 'optimized' to just use one `lfence`
before `rdtsc`:

``` .asm
// Linux kernel TSC usage (circa 2011+)
lfence
rdtsc
```

The kernel developer of the older `lfence` both sides approach appears to
object to this optimization as 'unsafe' due to being a barrier in only one
direction. The 'modern' thinking appears to be that while this is technically
true, a microprocessor would never take advantage of this reordering---there is
no performance reason to do so.

The Akaros project [investigated][arakos] a number of alternative approaches
(including all the above issues). Eventually taking the modern Linux approach
and suggesting the following:

``` .asm
// code before
lfence
rdtsc
// code to benchmark
lfence
rdtsc
// code after
```

Finally, for complete reference, an older [Intel guide][intel2] to using
`rdtsc` (pre `rdtscp` days) suggest that you 'warm up' the `cpuid` and `rdtsc`
instructions a few times before benchmarking the code:

``` .asm
// code before

// warmup
cpuid
rdtsc
cpuid
rdtsc
cpuid
rdtsc

cpuid
rdtsc
// code to benchmark
cpuid
rdtsc
// code afer
```

It's not clear if this is valuable any more when we have `rdtscp` to avoid
include the highly variable `cpuid` in our measurement region.

This is a very confusing situation. We keep it simple and stick with the
recommendation from Intel. This works well, but is a little more expensive due
to the `cpuid` calls compared to alternatives. However, it also appears to be
the 'safest', ensuring accurate measurements. For very frequent calls to the
TSC when benchmarking is not your goal, the standard Go `time.Now()` call is
very fast, essentially being `lfence; rdtsc`.

## Converting Cycles to Time

To convert from cycles to wall-clock time we need to know TSC frequency.
Frequency scaling on modern Intel chips doesn't affect the TSC.

Sadly, the only way to determine the TSC frequency appears to be through a MSR
using the `rdmsr` instruction. This instruction is privileged and can't be
executed from user-space.

If we could, we want to access the `MSR_PLATFORM_INFO`:

> Register Name: MSR_PLATFORM_INFO [15:8]
> Description: Package Maximum Non-Turbo Ratio (R/O)
>              The is the ratio of the frequency that invariant TSC runs at.
>              Frequency = ratio * 100 MHz.

The multiplicative factor of `100 MHz` varies across architectures. Luckily, it
appears to be `100 MHz` on all Intel architectures except Nehalem, for which it
is `133.3 MHz`.

If this method fails or is unavailable, Linux appears to determine the TSC
clock speed through a [calibration] [lxr3] against hardware timers.

For now, we don't provide the ability to convert cycles to time.

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

[intel1]: http://www.intel.com/content/www/us/en/embedded/training/ia-32-ia-64-benchmark-code-execution-paper.html
[intel2]: https://www.ccsl.carleton.ca/~jamuir/rdtscpm1.pdf
[lxr1]: http://lxr.free-electrons.com/source/include/asm-x86/system.h?v=2.6.25#L403
[lkml1]: https://lkml.org/lkml/2008/1/7/276
[lxr2]: http://lxr.free-electrons.com/source/arch/x86/include/asm/msr.h#L168
[lkml2]: https://lkml.org/lkml/2011/5/10/297
[arakos]: http://akaros.cs.berkeley.edu/lxr/akaros/kern/arch/x86/rdtsc_test.c
[lxr3]: http://lxr.free-electrons.com/source/arch/x86/kernel/tsc.c#L670

