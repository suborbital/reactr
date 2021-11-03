pub enum Method {
	GET,
	POST,
	PATCH,
	DELETE,
}

impl From<Method> for i32 {
	fn from(field_type: Method) -> Self {
		match field_type {
			Method::GET => 1,
			Method::POST => 2,
			Method::PATCH => 3,
			Method::DELETE => 4,
		}
	}
}
