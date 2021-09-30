package suborbital

// #include <stdint.h>
// void return_result(void* pointer, int32_t size, int32_t ident);
import "C"
import (
	"reflect"
	"runtime"
	"unsafe"
)

var RUNNABLE Runnable

// The Runnable interface is all that needs to be implemented by a Reactr runnable.
type Runnable interface {
	Run(input []byte) []byte
}

//export allocate
func allocate(size int32) uintptr {
	arr := make([]byte, size)

	header := (*reflect.SliceHeader)(unsafe.Pointer(&arr))

	runtime.KeepAlive(arr)

	return uintptr(header.Data)
}

//export deallocate
func deallocate(pointer uintptr, size int32) {
	var arr []byte

	header := (*reflect.SliceHeader)(unsafe.Pointer(&arr))
	header.Data = pointer
	header.Len = uintptr(size) // Interestingly, the types of .Len and .Cap here
	header.Cap = uintptr(size) // differ from standard Go, where they are both int

	arr = nil // I think this is sufficient to mark the slice for garbage collection
}

//export run_e
func run_e(pointer uintptr, size int32, ident int32) {
	var input []byte

	inputHeader := (*reflect.SliceHeader)(unsafe.Pointer(&input))
	inputHeader.Data = pointer
	inputHeader.Len = uintptr(size)
	inputHeader.Cap = uintptr(size)

	result := RUNNABLE.Run(input)

	runtime.KeepAlive(result)

	resultHeader := (*reflect.SliceHeader)(unsafe.Pointer(&result))

	C.return_result(unsafe.Pointer(resultHeader.Data), int32(len(result)), int32(ident))
}
