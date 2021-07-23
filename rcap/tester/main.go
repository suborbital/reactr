package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/suborbital/reactr/rcap"
)

func main() {
	gqlClient := rcap.DefaultGraphQLClient()

	resp, err := gqlClient.Do("https://api.rawkode.dev", `{
		allProfiles {
			forename
			surname
		}
	}
	`)

	if err != nil {
		fmt.Println(resp)
		log.Fatal(err)
	}

	jsonBytes, _ := json.Marshal(resp.Data)

	fmt.Println(string(jsonBytes))
}
