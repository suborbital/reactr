pub mod default_runnable;

use crate::{util, STATE};

use std::mem;
use std::slice;

extern "C" {
	fn return_result(result_pointer: *const u8, result_size: i32, ident: i32);
	fn return_error(code: i32, result_pointer: *const u8, result_size: i32, ident: i32);
}

pub struct RunErr {
	pub code: i32,
	pub message: String,
}

impl RunErr {
	pub fn new(code: i32, msg: &str) -> Self {
		RunErr {
			code,
			message: msg.into(),
		}
	}
}

pub struct HostErr {
	pub message: String,
}

impl HostErr {
	pub fn new(msg: &str) -> Self {
		HostErr {
			message: String::from(msg),
		}
	}
}

pub trait Runnable {
	fn run(&self, input: Vec<u8>) -> Result<Vec<u8>, RunErr>;
}

pub fn use_runnable(runnable: &'static dyn Runnable) {
	unsafe {
		STATE.runnable = runnable;
	}
}

/// # Safety
///
/// We hand over the the pointer to the allocated memory.
/// Caller has to ensure that the memory gets freed again.
#[no_mangle]
pub unsafe extern "C" fn allocate(size: i32) -> *const u8 {
	let mut buffer = Vec::with_capacity(size as usize);

	let pointer = buffer.as_mut_ptr();

	mem::forget(buffer);

	pointer as *const u8
}

/// # Safety
#[no_mangle]
pub unsafe extern "C" fn deallocate(pointer: *const u8, size: i32) {
	let _ = slice::from_raw_parts(pointer, size as usize);
}

/// # Safety
#[no_mangle]
pub unsafe extern "C" fn run_e(pointer: *const u8, size: i32, ident: i32) {
	STATE.ident = ident;

	// rebuild the memory into something usable
	let in_slice: &[u8] = slice::from_raw_parts(pointer, size as usize);

	let in_bytes = Vec::from(in_slice);

	let mut code = 0;

	// call the runnable and check its result
	let result: Vec<u8> = match STATE.runnable.run(in_bytes) {
		Ok(val) => val,
		Err(e) => {
			code = e.code;
			util::to_vec(e.message)
		}
	};

	let result_slice = result.as_slice();
	let result_size = result_slice.len();

	// call back to reactr to return the result or error
	if code != 0 {
		return_error(code, result_slice.as_ptr() as *const u8, result_size as i32, ident);
	} else {
		return_result(result_slice.as_ptr() as *const u8, result_size as i32, ident);
	}
}
