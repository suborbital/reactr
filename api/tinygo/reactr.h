#include <stdint.h>

typedef int32_t i32;

i32 get_ffi_result(void* ptr, i32 ident);
i32 add_ffi_var(void* name_ptr, i32 name_size, void* var_ptr, i32 val_size, i32 ident);

void return_result(void* rawdata, int32_t size, int32_t ident);
void return_error(int32_t code, void* rawdata, int32_t size, int32_t ident);

int32_t cache_set(void* key_ptr, int32_t key_size,  void* value_ptr, int32_t value_size, int32_t ttl, int32_t ident);
int32_t cache_get(void* key_ptr, int32_t key_size, int32_t ident);

void log_msg(void *ptr, i32 size, i32 level, i32 ident);

i32 request_get_field(i32 field_type, void* key_ptr, i32 key_size, i32 ident);
i32 request_set_field(i32 field_type, void *key_ptr, i32 key_size, void *value_ptr, i32 value_size, i32 ident);
