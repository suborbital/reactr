package runnable

// #include <stdint.h>
// int32_t cache_set(void* key_ptr, int32_t key_size,  void* value_ptr, int32_t value_size, int32_t ttl, int32_t ident);
// int32_t cache_get(void* key_ptr, int32_t key_size, int32_t ident);
import "C"

func CacheGet(key string) ([]byte, error) {
	ptr, size := rawSlicePointer([]byte(key))

	return result(C.cache_get(ptr, size, ident()))
}

func CacheSet(key, val string, ttl int) {
	keyPtr, keySize := rawSlicePointer([]byte(key))
	valPtr, valSize := rawSlicePointer([]byte(val))

	C.cache_set(keyPtr, keySize, valPtr, valSize, int32(ttl), ident())
	return
}
