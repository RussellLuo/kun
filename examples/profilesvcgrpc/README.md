# profilesvcgrpc

This example illustrates how to expose [profilesvc][1] as gRPC APIs.


## Prerequisites

- Protocol buffer compiler v3
- Go plugins for the protocol compiler

See [gRPC Go Quickstart][2] for installation instructions.


## Generate the code

```bash
$ go generate
```

## Test the server

Run the server:

```bash
$ go run cmd/main.go
2021/07/04 22:32:02 server listening at [::]:8080
```

Consume by [grpcurl][3]:

```bash
$ grpcurl -plaintext -d '{"profile": {"id": "1234", "name": "kok"}}' :8080 pb.Service/PostProfile
{

}
$ grpcurl -plaintext -d '{"id": "1234"}' :8080 pb.Service/GetProfile
{
  "profile": {
    "id": "1234",
    "name": "kok"
  }
}
```


[1]: https://github.com/RussellLuo/kok/tree/master/examples/profilesvc
[2]: http://www.grpc.io/docs/quickstart/go.html#prerequisites
[3]: https://github.com/fullstorydev/grpcurl