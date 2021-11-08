pub mod query;
use query::{QueryArg, QueryType};

use crate::{env, errors::HostResult, ffi};

// insert executes the pre-loaded database query with the name <name>,
// and passes the arguments defined by <args>
//
// the return value is the inserted auto-increment ID from the query result, if any,
// formatted as JSON with the key `lastInsertID`
pub fn insert(name: &str, args: Vec<QueryArg>) -> HostResult<Vec<u8>> {
	ffi::add_vars(args);
	// retreive the result from the host and return it
	ffi::result(env::db_exec(QueryType::INSERT.into(), name.as_ptr(), name.len() as i32))
}

// insert executes the pre-loaded database query with the name <name>,
// and passes the arguments defined by <args>
//
// the return value is the query result formatted as JSON, with each column name as a top-level key
pub fn select(name: &str, args: Vec<QueryArg>) -> HostResult<Vec<u8>> {
	ffi::add_vars(args);
	// retreive the result from the host and return it
	ffi::result(env::db_exec(QueryType::SELECT.into(), name.as_ptr(), name.len() as i32))
}
