use crate::ffi;
use crate::STATE;

pub struct QueryArg {
	pub name: String,
	pub value: String
}

extern {
	fn db_exec(query_type: i32, name_ptr: *const u8, name_size: i32, ident: i32) -> i32;
}

pub fn insert(name: &str, args: Vec<QueryArg>) -> Result<Vec<u8>, super::runnable::HostErr> {
	for a in args {
		super::ffi::add_var(a.name.as_str(), a.value.as_str())
	}

	let result_size = unsafe { db_exec(DB_QUERY_TYPE_INSERT, name.as_ptr(), name.len() as i32, super::STATE.ident) };

	// retreive the result from the host and return it
	super::ffi::result(result_size)
}

pub fn select(name: &str, args: Vec<QueryArg>) -> Result<Vec<u8>, super::runnable::HostErr> {
	for a in args {
		super::ffi::add_var(a.name.as_str(), a.value.as_str())
	}

	let result_size = unsafe { db_exec(DB_QUERY_TYPE_SELECT, name.as_ptr(), name.len() as i32, super::STATE.ident) };

	// retreive the result from the host and return it
	super::ffi::result(result_size)
}