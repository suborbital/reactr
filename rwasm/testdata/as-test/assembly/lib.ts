import { httpGet, logInfo } from "suborbital"

export function run(input: ArrayBuffer): ArrayBuffer {
	let inStr = String.UTF8.decode(input)
  
	let out = "hello, " + inStr

	logInfo(out)

	let resp = httpGet("https://1password.com")
	logInfo(String.UTF8.decode(resp))
  
	return String.UTF8.encode(out)
}