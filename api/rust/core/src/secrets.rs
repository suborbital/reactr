use crate::ffi;
use crate::runnable::HostErr;
use crate::STATE;
use crate::util;

extern "C" {
	fn get_secret_value(key_pointer: *const u8, key_size: i32, ident: i32) -> i32;
}

/// Fetches a secret value from the host
pub fn get_val(key: &str) -> Result<String, HostErr> {
	let result_size = unsafe { get_secret_value(key.as_ptr(), key.len() as i32, STATE.ident) };

	match ffi::result(result_size) {
		Ok(val) => return Ok(util::to_string(val)),
		Err(e) => return Err(e)
	}
}