use suborbital::runnable;
use suborbital::http;
use suborbital::util;
use suborbital::log;
use std::collections::BTreeMap;

struct Fetch{}

impl runnable::Runnable for Fetch {
    fn run(&self, input: Vec<u8>) -> Option<Vec<u8>> {
        let url = util::to_string(input);
        let result = http::get(url.as_str(), None);

        // test sending a POST request with headers and a body
        let mut headers = BTreeMap::new();
        headers.insert("Content-Type", "application/json");
        headers.insert("X-ATMO-TEST", "testvalgoeshere");

        let body = String::from("{\"message\": \"testing the echo!\"}").as_bytes().to_vec();

        let result2 = http::post("https://postman-echo.com/post", Some(body), Some(headers));
        log::info(util::to_string(result2).as_str());

        Some(result)
    }
}


// initialize the runner, do not edit below //
static RUNNABLE: &Fetch = &Fetch{};

#[no_mangle]
pub extern fn init() {
    runnable::set(RUNNABLE);
}
