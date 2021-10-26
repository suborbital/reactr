pub mod runnable;
pub mod db;
pub mod graphql;
pub mod http;
pub mod cache;
pub mod ffi;
pub mod file;
pub mod graphql;
pub mod http;
pub mod log;
pub mod req;
pub mod resp;
pub mod runnable;
pub mod util;

/// This file represents the Rust "API" for Reactr Wasm runnables. The functions defined herein are used to exchange
/// data between the host (Reactr, written in Go) and the Runnable (a Wasm module, in this case written in Rust).
use runnable::default_runnable::DefaultRunnable;

/// State struct to hold our dynamic Runnable
struct State<'a> {
	ident: i32,
	runnable: &'a dyn runnable::Runnable,
}

/// The state that holds the user-provided Runnable and the current ident
static mut STATE: State = State {
	ident: 0,
	runnable: &DefaultRunnable::new(),
};
