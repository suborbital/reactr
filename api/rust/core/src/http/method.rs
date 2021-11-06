pub enum Method {
	GET,
	HEAD,
	OPTIONS,
	POST,
	PUT,
	PATCH,
	DELETE,
}

impl From<Method> for i32 {
	fn from(field_type: Method) -> Self {
		match field_type {
			GET => 0,
			HEAD => 1,
			OPTIONS => 2,
			POST => 3,
			PUT => 4,
			PATCH => 5,
			DELETE => 6,
		}
	}
}
