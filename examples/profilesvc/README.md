# profilesvc

Let's take [profilesvc](https://github.com/go-kit/kit/tree/266ff8dc37c693d0649707e519c93c1f85868bdc/examples/profilesvc) as an example, see how we can generate the endpoint/http code based on the interface [Service](https://github.com/go-kit/kit/blob/266ff8dc37c693d0649707e519c93c1f85868bdc/examples/profilesvc/service.go#L9-L20).


## Prerequisites

1. Adjust the [Service](https://github.com/RussellLuo/kok/blob/master/examples/profilesvc/service.go#L11-L58) interface

    - Add a meaningful name to each input/output parameter, to get more human-readable field names in the corresponding request/response structs.
    - Add kok-specific comments (i.e. comments start with "// @kok") in a [OAS](http://spec.openapis.org/oas/v3.0.3)-inspired format, to describe the properties of the exposed HTTP APIs.

2. Implement `err2code()`

    - Provide a function named `err2code` in [err2code.go](err2code.go), which has the exactly same logic as [codeFrom](https://github.com/go-kit/kit/blob/266ff8dc37c693d0649707e519c93c1f85868bdc/examples/profilesvc/transport.go#L392-L401), to transform any business error to an HTTP code.


## Generate the code

1. Use the `kok` command

    Just run:

    ```bash
    $ kok ./service.go Service
    ```

2. Use `go:generate`

    Add `//go:generate kok ./service.go Service` before the Service interface, then run:

    ```bash
    $ go generate
    ```

## Test the service

Run the Profile service:

```bash
$ go run cmd/main.go
2020/05/08 10:22:22 transport=HTTP addr=:8080
```

Create a Profile:

```bash
$ curl -i -X POST http://localhost:8080/profiles -H "Content-Type: application/json" -d '{"profile": {"id": "1234", "name": "kok"}}'
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Fri, 08 May 2020 02:22:22 GMT
Content-Length: 3

{}
```

Get the profile you just created:

```bash
$ curl -i http://localhost:8080/profiles/1234
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Fri, 08 May 2020 02:22:25 GMT
Content-Length: 39

{"profile":{"id":"1234","name":"kok"}}
```
