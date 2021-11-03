pub enum Method {
	GET,
	POST,
	PATCH,
	DELETE,
}

impl From<Method> for i32 {
	use Method::*;
	fn from(field_type: Method) -> Self {
		match field_type {
			GET => 0,
			POST => 1,
			PUT => 2,
			PATCH => 3,
			DELETE => 4,
		}
	}
}
