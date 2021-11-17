//go:build tinygo.wasm

package runnable

// #include <reactr.h>
import "C"
import (
	"runtime"
)

func result(size int32) ([]byte, HostErr) {
	allocSize := size

	if size < 0 {
		if size == -1 {
			return nil, NewHostError("unknown error returned from host")
		}

		allocSize = -size
	}

	result := make([]byte, allocSize)
	resultPtr, _ := rawSlicePointer(result)

	if code := C.get_ffi_result(resultPtr, ident()); code != 0 {
		return nil, NewHostError("unknown error returned from host")
	}

	if size < 0 {
		return nil, NewHostError(string(result))
	}

	return result, nil
}

func AddVar(name, value string) {
	nameB := []byte(name)
	namePtr, nameSize := rawSlicePointer(nameB)
	runtime.KeepAlive(nameB)

	valueB := []byte(value)
	valuePtr, valueSize := rawSlicePointer(valueB)
	runtime.KeepAlive(valueB)

	C.add_ffi_var(namePtr, nameSize, valuePtr, valueSize, ident())
}
