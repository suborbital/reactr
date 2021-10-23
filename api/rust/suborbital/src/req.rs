pub mod field_type;

use crate::util;
use crate::ffi;
use crate::STATE;
use field_type::FieldType;

extern {
	fn request_get_field(field_type: i32, key_pointer: *const u8, key_size: i32, ident: i32) -> i32;
}

pub fn method() -> String {
	get_field(FieldType::Meta.into(), "method")
		.map_or("".into(), util::to_string)
}

pub fn url() -> String {
	get_field(FieldType::Meta.into(), "url")
		.map_or("".into(), util::to_string)
}

pub fn id() -> String {
	get_field(FieldType::Meta.into(), "id")
		.map_or("".into(), util::to_string)
}

pub fn body_raw() -> Vec<u8> {
	get_field(FieldType::Meta.into(), "body")
		.unwrap_or_default()
}

pub fn body_field(key: &str) -> String {
	get_field(FieldType::Body.into(), key)
		.map_or("".into(), util::to_string)
}

pub fn header(key: &str) -> String {
	get_field(FieldType::Header.into(), key)
		.map_or("".into(), util::to_string)
}

pub fn url_param(key: &str) -> String {
	get_field(FieldType::Params.into(), key)
		.map_or("".into(), util::to_string)
}

pub fn state(key: &str) -> Option<String> {
	get_field(FieldType::State.into(), key)
		.map(util::to_string)
}

pub fn state_raw(key: &str) -> Option<Vec<u8>> {
	get_field(FieldType::State.into(), key)
}

/// Executes the request via FFI
///
/// Then retreives the result from the host and returns it
fn get_field(field_type: i32, key: &str) -> Option<Vec<u8>> {
	let result_size = unsafe {
		request_get_field(field_type, key.as_ptr(), key.len() as i32, STATE.ident) 
	};

	ffi::result(result_size)
		.map_or(None, Option::from)
}
