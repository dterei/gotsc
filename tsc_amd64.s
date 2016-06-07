#include "textflag.h"

// func BenchStart() uint64
TEXT ·BenchStart(SB),NOSPLIT,$0-8
	CPUID
	RDTSC
	SHLQ	$32, DX
	ADDQ	DX, AX
	MOVQ	AX, ret+0(FP)
	RET

// func BenchEnd() uint64
TEXT ·BenchEnd(SB),NOSPLIT,$0-8
	BYTE	$0x0F // RDTSCP
	BYTE	$0x01
	BYTE	$0xF9
	SHLQ	$32, DX
	ADDQ	DX, AX
	MOVQ	AX, ret+0(FP)
	CPUID
	RET
