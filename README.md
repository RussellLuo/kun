# kok

The toolkit of [Go kit][1].


## Installation

```bash
$ go get -u github.com/RussellLuo/kok/cmd/kokgen
```

Usage:

```bash
$ kokgen -h
kokgen [flags] source-file interface-name
  -fmt
    	whether to make code formatted (default true)
  -out string
    	output directory (default ".")
  -pkg string
    	package name (default will infer)
  -test string
    	the YAML file that provides test-cases for HTTP (default "./http.test.yaml")
  -trace
    	whether to enable tracing
```


## HTTP API

### API annotations

1. Define the HTTP operation

    - Key: `@kok(op)`
    - Value: `"<method> <pattern>"`
        + **method**: The HTTP method
        + **pattern**: The URL pattern
    - Example:

        ```go
        type Service interface {
            // @kok(op): "POST /users"
            CreateUser(ctx context.Context) (err error)
        }
        ```

2. Define the HTTP request parameters

    - Key: `@kok(param)`
    - Value: `"name:<name>,type:<type>,in:<in>,alias:<alias>"`
        + **name**: The name of the argument in the interface method.
            - *Argument group*: By using `.` in **name**, multiple request arguments (each one is of basic type) can be grouped into one method argument (of struct type).
        + **type**: The type of the argument in the interface method.
            - Optional: Default will infer from the declaration of the interface method. **Required** for *argument group*.
        + **in**:
            - **path**: The argument is located in the request path.
            - **query**: The argument is located in the request query parameters.
            - **header**: The argument is located in the request header.
            - **cookie**: The argument is located in the request cookie.
                + Not supported
            - **body**: The argument is located in the request body.
                + Optional: All arguments, unless otherwise specified, are in body.
        + **alias**: The name of the argument in the request.
            - Optional: Defaults to **name** if not specified.
    - Example:
        + Simple argument:

            ```go
            type Service interface {
                // @kok(op): "DELETE /users/{id}"
                // @kok(param): "name:id,in:path"
                DeleteUser(ctx context.Context, id int) (err error)
            }

            // HTTP request: DELETE /users/101
            ```
        + Argument group:

            ```go
            type User struct {
                Name string
                Age  int
            }

            type Service interface {
                // @kok(op): "POST /users"
                // @kok(param): "name:user.Name,type:string,in:query,alias:name"
                // @kok(param): "name:user.Age,type:int,in:query,alias:age"
                CreateUser(ctx context.Context, user User) (err error)
            }

            // HTTP request: POST /users?name=russell&age=1
            ```

3. Define the status code of the success HTTP response

    - Key: `@kok(success)`
    - Value: `"statusCode:<statusCode>"`
        + **statusCode**: The status code of the success HTTP response.
            - Optional: Defaults to 200 if not specified.
    - Example:

        ```go
        type Service interface {
            // @kok(op): "POST /users"
            // @kok(success): "statusCode:201"
            CreateUser(ctx context.Context) (err error)
        }
        ```

### Encoding and decoding

See [HTTP codec interfaces](https://github.com/RussellLuo/kok/blob/master/pkg/codec/httpv2/codec.go).


## Examples

See [examples/profilesvc](examples/profilesvc).


## Documentation

Checkout the [Godoc][2].


## License

[MIT](LICENSE)


[1]: https://github.com/go-kit/kit
[2]: https://godoc.org/github.com/RussellLuo/kok
