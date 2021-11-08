use crate::env;

pub fn debug(msg: &str) {
	log_at_level(msg, 4)
}

pub fn info(msg: &str) {
	log_at_level(msg, 3)
}

pub fn warn(msg: &str) {
	log_at_level(msg, 2)
}

pub fn error(msg: &str) {
	log_at_level(msg, 1)
}

fn log_at_level(msg: &str, level: i32) {
	env::log_msg(msg.as_ptr(), msg.len() as i32, level)
}
