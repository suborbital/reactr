import { request_get_field } from "./env";
import { Result, ffi_result, getIdent, toFFI } from "./ffi"

const FIELD_TYPE_META: i32 = 0
const FIELD_TYPE_BODY: i32 = 1
const FIELD_TYPE_HEADER: i32 = 2
const FIELD_TYPE_PARAMS: i32 = 3
const FIELD_TYPE_STATE: i32 = 4

export function reqMethod(): Result {
	return get_field(FIELD_TYPE_META, "method")
}

export function reqURL(): Result {
	return get_field(FIELD_TYPE_META, "url")
}

export function reqID(): Result {
	return get_field(FIELD_TYPE_META, "id")
}

export function reqBody(): Result {
	return get_field(FIELD_TYPE_META, "body")
}

export function reqBodyField(key: string): Result {
	return get_field(FIELD_TYPE_BODY, key)
}

export function reqHeader(key: string): Result {
	return get_field(FIELD_TYPE_HEADER, key)
}

export function reqURLParam(key: string): Result {
	return get_field(FIELD_TYPE_PARAMS, key)
}

export function reqState(key: string): Result {
	return get_field(FIELD_TYPE_STATE, key)
}

export function reqStateRaw(key: string): Result {
	return get_field(FIELD_TYPE_STATE, key)
}

function get_field(field_type: i32, key: string): Result {
	let keyBuf = String.UTF8.encode(key)
	let keyFFI = toFFI(keyBuf)

	let result_size = request_get_field(field_type, keyFFI.ptr, keyFFI.size, getIdent())

	return ffi_result(result_size)
}