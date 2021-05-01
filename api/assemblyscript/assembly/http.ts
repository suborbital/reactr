import { fetch_url } from "./env"
import { ffi_result, getIdent, toFFI } from "./ffi"

export function httpGet(url: string): ArrayBuffer {
	return do_request(method_get, url, new ArrayBuffer(0))
}

export function httpPost(url: string, body: ArrayBuffer): ArrayBuffer {
	return do_request(method_post, url, body)
}

export function httpPatch(url: string, body: ArrayBuffer): ArrayBuffer {
	return do_request(method_patch, url, body)
}

export function httpDelete(url: string): ArrayBuffer {
	return do_request(method_delete, url, new ArrayBuffer(0))
}

const method_get = 1
const method_post = 2
const method_patch = 3
const method_delete = 4

function do_request(method: i32, url: string, body: ArrayBuffer): ArrayBuffer {
	// TODO: handle headers

	let urlBuf = String.UTF8.encode(url)
	let urlFFI = toFFI(urlBuf)

	let bodyFFI = toFFI(body)

	let result_size = fetch_url(method, urlFFI.ptr, urlFFI.size, bodyFFI.ptr, bodyFFI.size, getIdent())

	let result = ffi_result(result_size)

	return result
}