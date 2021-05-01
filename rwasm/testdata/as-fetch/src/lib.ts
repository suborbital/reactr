import { httpGet, logInfo } from "@suborbital/suborbital"

export function run(input: ArrayBuffer): ArrayBuffer {
	let url = String.UTF8.decode(input)

	logInfo("fetching " + url)

	let resp = httpGet(url)
  
	return resp
}