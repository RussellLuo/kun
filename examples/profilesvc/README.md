# profilesvc

Let's take [profilesvc](https://github.com/go-kit/kit/tree/266ff8dc37c693d0649707e519c93c1f85868bdc/examples/profilesvc) as an example, see how we can generate the endpoint/http code based on the interface [Service](https://github.com/go-kit/kit/blob/266ff8dc37c693d0649707e519c93c1f85868bdc/examples/profilesvc/service.go#L9-L20).


## Prerequisites

1. Adjust the [Service](https://github.com/RussellLuo/kok/blob/master/examples/profilesvc/service.go#L11-L39) interface

    - Add a meaningful name to each input/output parameter, to get more human-readable field names in the corresponding request/response structs.
    - Add kun-specific comments (i.e. comments start with "// @kun") in a [OAS](http://spec.openapis.org/oas/v3.0.3)-inspired format, to describe the properties of the exposed HTTP APIs.

2. Customize HTTP encoders and decoders

    - Override the method [Codec.EncodeFailureResponse](https://github.com/RussellLuo/kok/blob/master/examples/profilesvc/codec.go#L14-L18), to transform any business error to an HTTP response.

3. List business errors for generating failure responses in OAS (**Optional**)

    - See [GetFailures](https://github.com/RussellLuo/kok/blob/master/examples/profilesvc/codec.go#L37-L52).

4. Define HTTP test-cases in YAML (**Optional**)

    - See [http.test.yaml](http.test.yaml).


## Generate the code

1. Use the `kungen` command

    Just run:

    ```bash
    $ kungen ./service.go Service
    ```

2. Use `go:generate`

    Add `//go:generate kungen ./service.go Service` before the Service interface, then run:

    ```bash
    $ go generate
    ```

Code generated:

- [endpoint.go](endpoint.go)
- [http.go](http.go)
- [http_test.go](http_test.go)
- [http_client.go](http_client.go)
- [oasv2.go](oasv2.go)


## Test the server

Run the Profile server:

```bash
$ go run cmd/server/main.go
2020/08/05 20:37:36 transport=HTTP addr=:8080
```

Create a Profile:

```bash
$ curl -i -X POST http://localhost:8080/profiles -H "Content-Type: application/json" -d '{"profile": {"id": "1234", "name": "kun"}}'
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

{"profile":{"id":"1234","name":"kun"}}
```


## Generate OAS Documentation

```bash
$ curl http://localhost:8080/api > api.yaml
```

For those who want to preview the final generated documentation, see the pre-generated file [api.yaml](api.yaml).


## Test the client

Run the Profile client:

```bash
$ go run cmd/client/main.go
2020/08/05 20:37:45 GetProfile ok: {ID:1 Name:profile1 Addresses:[]}
2020/08/05 20:37:45 GetAddress ok: {ID:4 Location:address4}
2020/08/05 20:37:45 GetAddresses ok: [{ID:4 Location:address4}]
```


## Run tests:

```bash
$ go test -v -race
```
