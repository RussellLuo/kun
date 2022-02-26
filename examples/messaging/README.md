# messaging

This example illustrates how to specify multiple HTTP operations.


## Generate the code

```bash
$ go generate
```

## Test the server

Run the server:

```bash
$ go run cmd/main.go
2022/02/26 17:45:20 transport=HTTP addr=:8080
```

Get a message:

```bash
$ curl -XGET 'http://localhost:8080/messages/123456'
{"text":"user[]: message[123456]"}
```

Get a message from a specific user:

```bash
$ curl -XGET 'http://localhost:8080/users/me/messages/123456'
{"text":"user[me]: message[123456]"}
```
