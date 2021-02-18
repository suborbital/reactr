use suborbital::{http, runnable};

struct BenchFetch{}

impl runnable::Runnable for BenchFetch {
    fn run(&self, input: Vec<u8>) -> Option<Vec<u8>> {
        let url = String::from_utf8(input).unwrap();
        
        http::get(url.as_str(), None);

        Some("ok".as_bytes().to_vec())
    }
}


// initialize the runner, do not edit below //
static RUNNABLE: &BenchFetch = &BenchFetch{};

#[no_mangle]
pub extern fn init() {
    runnable::set(RUNNABLE);
}
