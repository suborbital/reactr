## Reactr FaaS ☁️

Reactr has built-in support for acting as a Functions-as-a-Service system (FaaS). Reactr Runnables can be exposed as endpoints to be triggered by an API request. Reactr will automatically create and manage the server on your behalf and allow for efficient execution of your Runnables over the network.

Reactr FaaS operates with similar semantics to the Reactr Go library. Jobs are triggered by making a POST request, and the caller receives a result ID in return. The job will be automatically scheduled and the result can be fetched later by passing the result ID in a later GET request. Results can optionally be waited for by passing the `?then=true` query parameter.

An example of creating a Reactr FaaS server can be found in [servertest](../rfaasservertest/main.go). Reactr FaaS uses Suborbital's [Vektor API framework](https://github.com/suborbital/vektor), and so all of its options are available, and the resulting server object can be optionally extended with other handlers. Below are the API calls available for Reactr jobs.

## Schedule a Job

URI: | `/do/:jobname`
:--- | :---
Method: | `POST`
Body: | Job payload (raw bytes)
Response: | JSON bytes representing the result
**Parameter** | **Effect**
 `then=true` | When provided, causes the request to wait until the scheduled job is completed, and returns the job result as raw bytes. If the job result was a struct, an attempt will be made to JSON marshal it before sending. If any error occurs, the response will have a non-200 HTTP status code and a body containing an error message.
 `callback={url}` | When provided, a webhook POST request will be sent to the provided URL when the job completes. The request will contain the bytes of the job result. If the job result was a struct, an attempt will be made to JSON marshal it before sending. If any error occurs, the request payload will be a string beginning with `job_err_result` followed by an error message. When `callback` is set, `then` will be ignored, and the response to the caller will be empty with HTTP status 200 OK.
**Example Request** | **Example Response**
`POST` `/do/compressimage` | `{"resultId":"7gj9n0adohm36zeqbfys4re6"}`

## Get a result

URI: | `/then/:resultid`
:--- | :---
Method: | `GET`
Body: | none
Response: | Job result (raw bytes)
**Example Request** | **Example Response**
`GET` `/then/7gj9n0adohm36zeqbfys4re6` | {job result bytes}
