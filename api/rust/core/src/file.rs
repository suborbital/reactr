use crate::{env, ffi};

/// Executes the request via FFI
///
/// Then retreives the result from the host and returns it
pub fn get_static(name: &str) -> Option<Vec<u8>> {
	let result_size = env::get_static_file(name.as_ptr(), name.len() as i32);

	match ffi::result(result_size) {
		Ok(res) => Some(res),
		Err(_) => None,
	}
}
