use suborbital::runnable::*;
use suborbital::log;
use suborbital::util;

use rustpython_vm as vm;
// this needs to be in scope in order to insert things into scope.globals
use vm::ItemProtocol;
use vm::builtins::{PyStrRef};

struct Python{}

impl Runnable for Python {
    fn run(&self, input: Vec<u8>) -> Result<Vec<u8>, RunErr> {
        let source = util::to_string(input);

        let vm_inst = vm::Interpreter::default();

        match vm_inst.enter(|vm| {
            run(source.as_str(), vm)
        }) {
            Ok(_) => log::info("all good"),
            Err(_) => {
                return Err(RunErr::new(500, "failed"))
            }
        };

        Ok(Vec::new())
    }
}

fn run(source: &str, vm: &vm::VirtualMachine) -> vm::PyResult<()> {
    let scope: vm::scope::Scope = vm.new_scope_with_builtins();
    
    scope
        .globals
        .set_item("log_info", vm.ctx.new_function("log_info", pylog).into(), vm)?;

    vm.compile(source, vm::compile::Mode::Eval, "<embedded>".to_owned())
        .map_err(|err| vm.new_syntax_error(&err))
        .and_then(|code_obj| {
            vm.run_code_obj(code_obj, scope.clone())
        })
        .map_err(|err| {
            vm.print_exception(err.clone());
            err
        })?;

    Ok(())
}

fn pylog(msg: PyStrRef) {
    log::info(msg.as_str());
}

// initialize the runner, do not edit below //
static RUNNABLE: &Python = &Python{};

#[no_mangle]
pub extern fn _start() {
    use_runnable(RUNNABLE);
}
