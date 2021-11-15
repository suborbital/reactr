package runnable

// #include <stdint.h>
// int32_t get_ffi_result(void* ptr, int32_t ident);
// int32_t add_ffi_var(void* name_ptr, int32_t name_size, void* var_ptr, int32_t val_size, int32_t ident);
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
