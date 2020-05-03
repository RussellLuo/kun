# profilesvc

Let's take [profilesvc](https://github.com/go-kit/kit/tree/266ff8dc37c693d0649707e519c93c1f85868bdc/examples/profilesvc) as an example, see how we can generate endpoint/http code based on the interface [Service](https://github.com/go-kit/kit/blob/266ff8dc37c693d0649707e519c93c1f85868bdc/examples/profilesvc/service.go#L9-L20).


## Adjust the Service interface

1. Add a meaningful name to all input/output parameters, for human-readable field names in the corresponding request/response struct.
2. Add kok-specific comments (i.e. comments start with "// @kok") in a OpenAPI-inspired format, to describe the HTTP property.


## Implement `err2code()`

Provide a function named `err2code` in [err2code.go](err2code.go), which is the same as [codeFrom](https://github.com/go-kit/kit/blob/266ff8dc37c693d0649707e519c93c1f85868bdc/examples/profilesvc/transport.go#L392-L401).


## Generate code

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
2020/05/03 18:02:52 transport=HTTP addr=:8080
```

Create a Profile:

```bash
$ curl -X POST http://localhost:8080/profiles -H "Content-Type: application/json" -d '{"profile": {"id":"1234","name":"Go Kit"}}'
{}
```

Get the profile you just created:

```bash
$ curl http://localhost:8080/profiles/1234
{"profile":{"id":"1234","name":"Go Kit"}}
```
