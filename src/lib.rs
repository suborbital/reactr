use std::ffi::{CStr, CString};
use std::mem;
use std::os::raw::{c_char, c_void};
use wasm_bindgen::prelude::*;

mod run;

#[no_mangle]
#[wasm_bindgen]
pub extern fn allocate_input(size: usize) -> *mut c_void {
    let mut buffer = Vec::with_capacity(size);
    let pointer = buffer.as_mut_ptr();
    mem::forget(buffer);

    pointer as *mut c_void
}

#[no_mangle]
#[wasm_bindgen]
pub extern fn deallocate(pointer: *mut c_void, capacity: usize) {
    unsafe {
        let _ = Vec::from_raw_parts(pointer, 0, capacity);
    }
}

#[no_mangle]
#[wasm_bindgen]
pub extern fn run_e(input: *mut c_char) -> *mut c_char {
    // yes this is a convoluted mess, but it's what's needed

    // convert pointer to cstring
    let in_cstr = unsafe { CStr::from_ptr(input) };

    // convert cstring to String
    let in_str = String::from(in_cstr.to_str().unwrap());

    // run the runnable and convert it to a Vec
    let output = run::run(in_str).unwrap().as_bytes().to_vec();

    // convert the Vec to a pointer and return
    let owned = unsafe { CString::from_vec_unchecked(output) };
    return owned.into_raw();
}