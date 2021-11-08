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

pub fn runnable_unset() -> RunErr {
	RunErr::new(-1, "No runnable set")
}

pub type HostResult<T> = Result<T, HostErr>;
pub type RunResult<T> = Result<T, RunErr>;
