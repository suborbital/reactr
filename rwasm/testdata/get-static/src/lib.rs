use suborbital::runnable;
use suborbital::file;

struct GetStatic{}

impl runnable::Runnable for GetStatic {
    fn run(&self, input: Vec<u8>) -> Option<Vec<u8>> {
        let in_string = String::from_utf8(input).unwrap();
    
        let file = file::get_static(in_string.as_str())
            .unwrap_or("".as_bytes().to_vec());

        Some(file)
    }
}


// initialize the runner, do not edit below //
static RUNNABLE: &GetStatic = &GetStatic{};

#[no_mangle]
pub extern fn init() {
    runnable::set(RUNNABLE);
}
