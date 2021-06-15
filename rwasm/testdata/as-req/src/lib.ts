import { logInfo, reqBody, reqBodyField, reqID, reqMethod, reqState, reqURL } from "@suborbital/suborbital"

export function run(input: ArrayBuffer): ArrayBuffer {
	let inStr = String.UTF8.decode(input)
	logInfo(inStr)
  
	logInfo(reqMethod())
	logInfo(String.UTF8.decode(reqBody()))
	logInfo(reqBodyField("username"))
	logInfo(reqBodyField("baz")) // ensure it doesn't crash on something that doesn't exist
	logInfo(reqURL())
	logInfo(reqID())

	let hello = reqState("hello")
	logInfo(hello)

	return String.UTF8.encode("hello " + hello)
}