import { fetch_url } from "./env"
import { Result, ffi_result, getIdent, toFFI } from "./ffi"

export function httpGet(url: string, headers: Map<string, string> | null): Result {
	return do_request(method_get, url, new ArrayBuffer(0), headers)
}

export function httpHead(url: string, headers: Map<string, string> | null): Result {
	return do_request(method_head, url, new ArrayBuffer(0), headers)
}

export function httpOptions(url: string, headers: Map<string, string> | null): Result {
	return do_request(method_options, url, new ArrayBuffer(0), headers)
}

export function httpPost(url: string, body: ArrayBuffer, headers: Map<string, string> | null): Result {
	return do_request(method_post, url, body, headers)
}

export function httpPut(url: string, body: ArrayBuffer, headers: Map<string, string> | null): Result {
	return do_request(method_put, url, body, headers)
}

export function httpPatch(url: string, body: ArrayBuffer, headers: Map<string, string> | null): Result {
	return do_request(method_patch, url, body, headers)
}

export function httpDelete(url: string, headers: Map<string, string> | null): Result {
	return do_request(method_delete, url, new ArrayBuffer(0), headers)
}

const method_get = 0
const method_head = 1
const method_options = 2
const method_post = 3
const method_put = 4
const method_patch = 5
const method_delete = 6

function do_request(method: i32, url: string, body: ArrayBuffer, headers: Map<string, string> | null): Result {
	var headerString = ""
	if (headers != null) {
		headerString = renderHeaderString(headers)
	}

	let urlBuf = String.UTF8.encode(url + headerString)
	let urlFFI = toFFI(urlBuf)

	let bodyFFI = toFFI(body)

	let result_size = fetch_url(method, urlFFI.ptr, urlFFI.size, bodyFFI.ptr, bodyFFI.size, getIdent())

	return ffi_result(result_size)
}

function renderHeaderString(headers: Map<string,string>): string {
	var rendered: string = ""
	let keys = headers.keys()
	
	for (let i = 0; i < keys.length; ++i) {
		let key = keys[i]
		let val = headers.get(key)

		rendered += "::" + key + ":" + val
	}

	return rendered
}