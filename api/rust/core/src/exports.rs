use std::mem;

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
	crate::runnable::set_ident(ident);

	// rebuild the memory into something usable
	let in_bytes = Vec::from_raw_parts(pointer, size as usize, size as usize);
	crate::runnable::execute(in_bytes);
}
