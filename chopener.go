package chopener

import (
	"unsafe"
)

type hchan struct {
	qcount   uint
	dataqsiz uint
	buf      unsafe.Pointer
	elemsize uint16
	closed   uint32
}

//go:noinline
func Open[T any](ch *chan T) {
	(*hchan)(unsafe.Pointer(*(*uintptr)(unsafe.Pointer(ch)))).closed = 0
}
