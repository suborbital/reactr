use suborbital::runnable::*;
use suborbital::http;
use suborbital::util;
use suborbital::log;
use std::collections::BTreeMap;

struct Fetch{}

impl Runnable for Fetch {
    fn run(&self, input: Vec<u8>) -> Result<Vec<u8>, RunErr> {
        let url = util::to_string(input);
        http::get(url.as_str(), None)?;

        // test sending a POST request with headers and a body
        let mut headers = BTreeMap::new();
        headers.insert("Content-Type", "application/json");
        headers.insert("X-ATMO-TEST", "testvalgoeshere");

        let body = String::from("{\"message\": \"testing the echo!\"}").as_bytes().to_vec();

        match http::post("https://postman-echo.com/post", Some(body), Some(headers)) {
            Ok(res) => {
                log::info(util::to_string(res.clone()).as_str());
                Ok(res)
            },
            Err(e) => Err(e)
        }
    }
}


// initialize the runner, do not edit below //
static RUNNABLE: &Fetch = &Fetch{};

#[no_mangle]
pub extern fn init() {
    use_runnable(RUNNABLE);
}
