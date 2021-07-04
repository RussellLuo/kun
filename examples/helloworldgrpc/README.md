# helloworldgrpc

This example illustrates how to expose [helloworld][1] as gRPC APIs.


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
2021/07/04 18:18:00 server listening at [::]:8080
```

Consume by [grpcurl][3]:

```bash
$ grpcurl -plaintext -d '{"name": "Tracey"}' :8080 pb.Service/SayHello
{
  "message": "Hello Tracey"
}
```


[1]: https://github.com/RussellLuo/kok/tree/master/examples/helloworld
[2]: http://www.grpc.io/docs/quickstart/go.html#prerequisites
[3]: https://github.com/fullstorydev/grpcurl