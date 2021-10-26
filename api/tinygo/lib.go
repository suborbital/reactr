package runnable

// #include <stdint.h>
// void return_result(void* rawdata, int32_t size, int32_t ident);
// void return_error(int32_t code, void* rawdata, int32_t size, int32_t ident);
import "C"
import (
	"errors"
	"reflect"
	"runtime"
	"unsafe"
)

var runnable_ Runnable

func Use(runnable Runnable) {
	runnable_ = runnable
}

// The Runnable interface is all that needs to be implemented by a Reactr runnable.
type Runnable interface {
	Run(input []byte) ([]byte, error)
}

// RunErr adds a status code for use in FFI calls to return_result()
type RunErr struct {
	error
	Code int
}

// NewError creates a new RunErr, which is a normal Go error plus an error code
func NewError(code int, message string) RunErr {
	return RunErr{errors.New(message), code}
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
func run_e(rawdata uintptr, size int32, ident int32) {
	var input []byte

	inputHeader := (*reflect.SliceHeader)(unsafe.Pointer(&input))
	inputHeader.Data = rawdata
	inputHeader.Len = uintptr(size)
	inputHeader.Cap = uintptr(size)

	result, err := runnable_.Run(input)

	if err != nil {
		returnError(err, ident)
		return
	}

	resPtr, resLen := rawSlicePointer(result)

	C.return_result(resPtr, resLen, ident)
}

func returnError(err error, ident int32) {
	code := int32(500)

	if err == nil {
		C.return_error(code, unsafe.Pointer(uintptr(0)), 0, ident)
		return
	}

	switch e := err.(type) {
	case RunErr:
		code = int32(e.Code)
	}

	errPtr, errLen := rawSlicePointer([]byte(err.Error()))

	C.return_error(code, errPtr, errLen, ident)
}

func rawSlicePointer(slice []byte) (unsafe.Pointer, int32) {
	header := (*reflect.SliceHeader)(unsafe.Pointer(&slice))

	return unsafe.Pointer(header.Data), int32(len(slice))
}
