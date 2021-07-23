import { graphQLQuery, logInfo } from "@suborbital/suborbital"

export function run(_: ArrayBuffer): ArrayBuffer {
	let result = graphQLQuery("https://api.rawkode.dev", "{ allProfiles { forename, surname } }")
	if (result.byteLength == 0) {
		return String.UTF8.encode("failed")
	}

	logInfo(String.UTF8.decode(result))

	return result
}