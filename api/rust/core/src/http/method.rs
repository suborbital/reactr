pub enum Method {
	GET,
	POST,
	PATCH,
	DELETE,
}

impl From<Method> for i32 {
	fn from(field_type: Method) -> Self {
		match field_type {
			Method::GET => 0,
			Method::POST => 1,
			Method::PATCH => 2,
			Method::DELETE => 3,
		}
	}
}
