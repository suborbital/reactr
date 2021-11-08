#![allow(clippy::not_unsafe_ptr_arg_deref)]
suborbital_macro::wrap_host_functions! {
extern "C" {
	pub fn cache_set(
		key_pointer: *const u8,
		key_size: i32,
		value_pointer: *const u8,
		value_size: i32,
		ttl: i32,
		ident: i32,
	) -> i32;
	pub fn cache_get(key_pointer: *const u8, key_size: i32, ident: i32) -> i32;

	// database
	pub fn db_exec(query_type: i32, name_ptr: *const u8, name_size: i32, ident: i32) -> i32;

	// FFI
	pub fn get_ffi_result(pointer: *const u8, ident: i32) -> i32;
	pub fn add_ffi_var(name_ptr: *const u8, name_len: i32, val_ptr: *const u8, val_len: i32, ident: i32) -> i32;

	// file
	pub fn get_static_file(name_ptr: *const u8, name_size: i32, ident: i32) -> i32;

	// graphql
	pub fn graphql_query(
		endpoint_pointer: *const u8,
		endpoint_size: i32,
		query_pointer: *const u8,
		query_size: i32,
		ident: i32,
	) -> i32;

	// https
	pub fn fetch_url(
		method: i32,
		url_pointer: *const u8,
		url_size: i32,
		body_pointer: *const u8,
		body_size: i32,
		ident: i32,
	) -> i32;
	// Log

	pub fn log_msg(pointer: *const u8, result_size: i32, level: i32, ident: i32);

	/// Return with result
	pub fn return_result(result_pointer: *const u8, result_size: i32, ident: i32);

	// Return with Run Error
	pub fn return_error(code: i32, result_pointer: *const u8, result_size: i32, ident: i32);

	// Request

	pub fn request_get_field(field_type: i32, key_pointer: *const u8, key_size: i32, ident: i32) -> i32;
	pub fn request_set_field(
		field_type: i32,
		key_pointer: *const u8,
		key_size: i32,
		val_pointer: *const u8,
		val_size: i32,
		ident: i32,
	) -> i32;

	// Response

	pub fn resp_set_header(key_pointer: *const u8, key_size: i32, val_pointer: *const u8, val_size: i32, ident: i32);

  }
}
