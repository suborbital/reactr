use suborbital::runnable;
use suborbital::log;
use suborbital::req;
use suborbital::util;

struct Log{}

impl runnable::Runnable for Log {
    fn run(&self, input: Vec<u8>) -> Option<Vec<u8>> {
        let in_string = String::from_utf8(input).unwrap();
        log::info(in_string.as_str());

        log::info(req::method().as_str());
        log::info(util::to_string(req::body_raw()).as_str());
        log::info(req::body_field("username").as_str());
        log::info(req::body_field("baz").as_str()); // testing for a field that doesn't exist
        log::info(req::url().as_str());
        log::info(req::id().as_str());
        log::info(req::state("hello").as_str());
    
        Some(String::from(format!("hello {}", req::state("hello"))).as_bytes().to_vec())
    }
}


// initialize the runner, do not edit below //
static RUNNABLE: &Log = &Log{};

#[no_mangle]
pub extern fn init() {
    runnable::set(RUNNABLE);
}
