use crate::STATE;
use std::mem;

use crate::env;
use crate::errors::RunErr;

pub fn return_error(err: RunErr) {
	let RunErr { code, message } = err;
	env::return_error(code, message.as_ptr(), message.len() as i32)
}

pub fn return_result(result_data: Vec<u8>) {
	env::return_result(result_data.as_ptr(), result_data.len() as i32)
}

pub trait Runnable {
	fn run(&self, input: Vec<u8>) -> Result<Vec<u8>, RunErr>;
}

pub fn use_runnable(runnable: &'static dyn Runnable) {
	unsafe {
		STATE.runnable = Some(runnable);
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
pub unsafe extern "C" fn deallocate(pointer: *mut u8, size: i32) {
	drop(Vec::from_raw_parts(pointer, size as usize, size as usize))
}

/// # Safety
#[no_mangle]
pub unsafe extern "C" fn run_e(pointer: *mut u8, size: i32, ident: i32) {
	STATE.ident = ident;

	// rebuild the memory into something usable
	let in_bytes = Vec::from_raw_parts(pointer, size as usize, size as usize);

	match execute_runnable(STATE.runnable, in_bytes) {
		Ok(data) => {
			return_result(data);
		}
		Err(err) => {
			return_error(err);
		}
	}
}

fn execute_runnable(runnable: Option<&dyn Runnable>, data: Vec<u8>) -> Result<Vec<u8>, RunErr> {
	if let Some(runnable) = runnable {
		return runnable.run(data);
	}
	Err(RunErr::new(-1, "No runnable set"))
}
