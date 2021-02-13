#ifndef runnable_api_h
#define runnable_api_h

#include <stdlib.h>
#include <stdint.h>

#if __wasm32__

__attribute__((__import_name__("return_result")))
extern void return_result(char *result_pointer, int result_size, int ident);

#endif
#endif /* runnable_api_h */
