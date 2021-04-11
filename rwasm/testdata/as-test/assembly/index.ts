// DO NOT EDIT; generated file

import { return_result, fromFFI, getIdent, setIdent } from "suborbital";
import { run } from "./lib"

export function run_e(ptr: usize, size: i32, ident: i32): void {
  // set the current ident for other API methods to use
	setIdent(ident)

  // read the memory that was passed as input
	var inBuffer = fromFFI(ptr, size)

  // execute the Runnable
	let result = run(inBuffer)

  // return the result to the host
  return_result(changetype<usize>(result), result.byteLength, getIdent())
}

export function allocate(size: i32): usize {
  return heap.alloc(size)
}

export function deallocate(ptr: i32, _: i32): void {
  heap.free(ptr)
}