
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