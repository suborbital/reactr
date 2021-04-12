// returns a result to the host
export declare function return_result(ptr: usize, size: i32, ident: i32): void
// logs a message using the hosts' logger
export declare function log_msg(ptr: usize, size: i32, level: i32, ident: i32): void
// makes an http request
export declare function fetch_url(method: i32, url_ptr: usize, url_size: i32, body_ptr: usize, body_size: i32, ident: i32): i32
// gets the result of a guest->host FFI call
export declare function get_ffi_result(ptr: usize, ident: i32): i32