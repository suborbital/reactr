use suborbital::runnable;
use suborbital::cache;

struct RustSet{}

impl runnable::Runnable for RustSet {
    fn run(&self, input: Vec<u8>) -> Option<Vec<u8>> {
        cache::set("important", input, 0);
    
        Some(String::from("hello").as_bytes().to_vec())
    }
}


// initialize the runner, do not edit below //
static RUNNABLE: &RustSet = &RustSet{};

#[no_mangle]
pub extern fn init() {
    runnable::set(RUNNABLE);
}
