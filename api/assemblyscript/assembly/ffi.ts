import { get_ffi_result } from "./env"
import { fromFFI } from "./util"

var current_ident: i32 = 0;

export function setIdent(ident: i32): void {
	current_ident = ident
}

export function getIdent(): i32 {
	return current_ident
}

export function ffi_result(size: i32): ArrayBuffer {
	if (size < 0) {
		return new ArrayBuffer(0)
	}

	let result_ptr = heap.alloc(size)

	let code = get_ffi_result(result_ptr, current_ident)
	if (code != 0) {
		return new ArrayBuffer(0)
	}

	return fromFFI(result_ptr, size)
}