# usersvc

This example illustrates how to apply *Argument aggregation* in struct fields.


## Generate the code

```bash
$ go generate
```

## Test the server

Run the server:

```bash
$ go run cmd/main.go
2021/11/14 18:13:15 transport=HTTP addr=:8080
```

Create a user:

```bash
$ curl -XPOST 'http://localhost:8080/users?name=Tracey&age=1'
{"Name":"Tracey","Age":1,"IP":"127.0.0.1"}
```

Create a user with the `X-Forwarded-For` header:

```bash
$ curl -H 'X-Forwarded-For: 192.168.0.1' -XPOST 'http://localhost:8080/users?name=Tracey&age=1'
{"Name":"Tracey","Age":1,"IP":"192.168.0.1"}
```
