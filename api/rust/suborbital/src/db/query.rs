pub enum QueryType {
	SELECT,
	INSERT,
}

impl From<QueryType> for i32 {
	fn from(query_type: QueryType) -> Self {
		match query_type {
			QueryType::INSERT => 0,
			QueryType::SELECT => 1,
		}
	}
}