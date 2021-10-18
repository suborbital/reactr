use suborbital::runnable::*;
use suborbital::db;

struct RsDbtest{}

impl Runnable for RsDbtest {
    fn run(&self, _: Vec<u8>) -> Result<Vec<u8>, RunErr> {
        let mut args: Vec<db::QueryArg> = Vec::new();
        // args.push(db::QueryArg{name: String::from("uuid"), value: String::from("qwertyuiop")});
        args.push(db::QueryArg{name: String::from("email"), value: String::from("connor@suborbital.dev")});

        match db::select("SelectUserWithEmail", args) {
            Ok(result) => Ok(result),
            Err(e) => {
                Err(RunErr::new(500, e.message.as_str()))
            }
        }
    }
}


// initialize the runner, do not edit below //
static RUNNABLE: &RsDbtest = &RsDbtest{};

#[no_mangle]
pub extern fn _start() {
    use_runnable(RUNNABLE);
}
