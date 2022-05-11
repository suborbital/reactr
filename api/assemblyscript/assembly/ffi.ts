import { get_ffi_result } from "./env"

var current_ident: i32 = 0;

export class Result {
	Result: ArrayBuffer
	Err: Error | null

	constructor(result: ArrayBuffer, err: Error | null) {
		this.Result = result
		this.Err = err
	}

	toString(): string {
		let err = this.Err
		if(err) {
			return err.toString()
		}

		return String.UTF8.decode(this.Result)
	}
}

export function setIdent(ident: i32): void {
	current_ident = ident
}

export function getIdent(): i32 {
	return current_ident
}

export function ffi_result(size: i32): Result {
	let allocSize = size

	let unknownRes: Result = new Result(new ArrayBuffer(0), new Error("unknown error returned from host"))

	if (size < 0) {
		if (size == -1) {
			return unknownRes
		}

		allocSize = size * -1
	}

	let result_ptr = heap.alloc(allocSize)

	let code = get_ffi_result(result_ptr, current_ident)
	if (code != 0) {
		heap.free(result_ptr)
		return unknownRes
	}

	let data = fromFFI(result_ptr, allocSize)
	heap.free(result_ptr)

	if (size < 0) {
		return new Result(new ArrayBuffer(0), new Error(String.UTF8.decode(data)))
	}

	return new Result(data, null)
}

export function fromFFI(ptr: usize, len: i32): ArrayBuffer {
	let mem = new Uint8Array(len)

	for (let i = 0; i < len; i++) {
		mem[i] = load<u8>(ptr + i);
	}

	return mem.buffer
}

export class ffiValue {
	ptr: usize
	size: i32

	constructor(ptr: usize, size: i32) {
		this.ptr = ptr
		this.size = size
	}
}

export function toFFI(val: ArrayBuffer): ffiValue {
	return new ffiValue(changetype<usize>(val), val.byteLength)
}