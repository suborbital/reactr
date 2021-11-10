use once_cell::unsync::OnceCell;

use crate::env;
use crate::errors::{runnable_unset, RunErr};

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

pub(crate) fn execute(input: Vec<u8>) {
	unsafe {
		STATE
			.get()
			.map_or_else(|| Err(runnable_unset()), |state| state.runnable.run(input))
			.map_or_else(return_error, return_result)
	}
}

pub fn set_runnable(runnable: &'static dyn Runnable) {
	unsafe {
		if STATE.set(State { runnable, ident: 0 }).is_err() {
			return_error(RunErr::new(0, "Can only set runable once"))
		}
	}
}

/// This file represents the Rust "API" for Reactr Wasm runnables. The functions defined herein are used to exchange
/// data between the host (Reactr, written in Go) and the Runnable (a Wasm module, in this case written in Rust).

/// State struct to hold our dynamic Runnable
pub(crate) struct State<'a> {
	ident: i32,
	runnable: &'a dyn Runnable,
}

/// The state that holds the user-provided Runnable and the current ident
static mut STATE: OnceCell<State> = OnceCell::new();

pub(crate) fn current_ident() -> i32 {
	unsafe { STATE.get().unwrap().ident }
}

pub(crate) fn set_ident(i: i32) {
	unsafe {
		STATE
			.get_mut()
			.map_or_else(|| return_error(runnable_unset()), |state| state.ident = i)
	}
}
