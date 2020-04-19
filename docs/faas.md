## Hive FaaS

Hive has built-in support for acting as a Functions-as-a-Service system (FaaS). Hive Runnables can be exposed as endpoints to be triggered by an API request. Hive will automatically create and manage the server on your behalf and allow for efficient execution of your Runnables over the network.

Hive FaaS operates with similar semantics to the Hive Go library. Jobs are triggered by making a POST request, and the caller receives a result ID in return. The job will be automatically scheduled and the result can be fetched later by passing the result ID in a later GET request. Results can optionally be waited for by passing the `?then=true` query parameter.

An example of creating a Hive FaaS server can be found in [servertest](../servertest/main.go). Hive FaaS uses SubOrbital's [GusT API framework](https://github.com/suborbital/gust), and so all of its options are available, and the resulting server object can be optionally extended with other handlers. Below are the API calls available for Hive jobs.

## Schedule a Job

URI: | `/do/:jobname`
--- | ---
Method: | `POST`
Body: | Job payload (raw bytes)
Response: | JSON bytes representing the result
**Parameter** | **Effect**
 `?then=true` | when provided, causes the request to wait until the scheduled job is completed, and returns the job result as raw bytes
**Example Request** | **Example Response**
`POST` `/do/compressimage` | `{"resultId":"7gj9n0adohm36zeqbfys4re6"}`

## Get a result

URI: | `/then/:resultid`
--- | ---
Method: | `GET`
Body: | none
Response: | Job result (raw bytes)
**Example Request** | **Example Response**
`GET` `/then/7gj9n0adohm36zeqbfys4re6` | {job result bytes}
