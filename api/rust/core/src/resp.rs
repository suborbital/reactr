use crate::env;

pub fn set_header(key: &str, val: &str) {
	env::resp_set_header(key.as_ptr(), key.len() as i32, val.as_ptr(), val.len() as i32);
}

pub fn content_type(ctype: &str) {
	set_header("Content-Type", ctype);
}
