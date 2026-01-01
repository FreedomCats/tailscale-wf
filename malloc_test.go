// Copyright (c) 2021 The Inet.Af AUTHORS. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package wf

import (
	"testing"
	"unsafe"
)

// Odd-sized allocations (UTF-16 strings, byte blobs such as app IDs) must
// not misalign the ones that follow: structs with pointers are copied into
// arena memory, and a misaligned destination crashes the process with
// "bulkBarrierPreWrite: unaligned arguments" whenever the GC is marking.
func TestArenaAllocAligned(t *testing.T) {
	var a arena
	defer a.Dispose()

	align := unsafe.Alignof(uintptr(0))
	for _, n := range []uintptr{1, 3, 7, 8, 13, 2, 41, 6, 5, 4089, 1, 9} {
		p := uintptr(a.Alloc(n))
		if p%align != 0 {
			t.Errorf("Alloc(%d) = %#x, not aligned to %d", n, p, align)
		}
	}
}

// Rounding up must not break a full-slab allocation.
func TestArenaAllocFullSlab(t *testing.T) {
	var a arena
	defer a.Dispose()

	a.Alloc(1)
	if p := a.Alloc(slabSize); p == nil {
		t.Fatal("Alloc(slabSize) returned nil")
	}
	if len(a.slabs) != 2 {
		t.Errorf("slabs = %d, want 2", len(a.slabs))
	}
}
