# kok

The toolkit of [Go kit][1].


## Features

1. Code Generation Tool

    - [x] HTTP
        + [x] HTTP Server
        + [x] HTTP Tests
        + [x] HTTP Client
        + [x] OAS Documentation
    - [ ] gRPC

2. Useful Packages

    - [appx](pkg/appx): Application framework for HTTP and CRON applications.
    - [prometheus](pkg/prometheus): Prometheus metrics utilities.
    - [werror](pkg/werror): Classified business errors.


## Installation

```bash
$ go get -u github.com/RussellLuo/kok/cmd/kokgen
```

<details>
  <summary> Usage </summary>

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

</details>


## HTTP API

### API annotations

<details>
  <summary> Define the HTTP request operation </summary>

- Key: `@kok(op)`
- Value: `"<method> <pattern>"`
    + **method**: The request method
    + **pattern**: The request URL
- Example:

    ```go
    type Service interface {
        // @kok(op): "POST /users"
        CreateUser(ctx context.Context) (err error)
    }
    ```

</details>

<details>
  <summary> Define the HTTP request parameters </summary>

- Key: `@kok(param)`
- Value: `"name:<name>,type:<type>,in:<in>,alias:<alias>"`
    + **name**: The name of the method argument.
        - *Argument group*: By using `.` in **name**, multiple request parameters (each one is of basic type) can be grouped into one method argument (of struct type).
    + **type**: The type of the method argument.
        - Optional: Default will infer from the method declaration.
        - **Required** for arguments in *Argument group*.
    + **in**:
        - **path**: The method argument is passed via the request path.
        - **query**: The method argument is passed via the request query string.
        - **header**: The method argument is passed via the request headers.
        - **cookie**: The method argument is passed via the request cookies.
            + Not supported yet
        - **body**: The method argument is passed via the request body.
            + Optional: All method arguments, unless otherwise specified, are in **body**.
    + **alias**: The name of the request parameter.
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

        // HTTP request: POST /users?name=tracey&age=1
        ```

</details>

<details>
  <summary> Define the status code of the success HTTP response </summary>


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

</details>

### Encoding and decoding

See [HTTP codec interfaces](https://github.com/RussellLuo/kok/blob/master/pkg/codec/httpv2/codec.go).


## Examples

See [examples/profilesvc](examples/profilesvc).


## Documentation

Checkout the [Godoc][3].


## License

[MIT](LICENSE)


[1]: https://github.com/go-kit/kit
[2]: https://github.com/RussellLuo/appx
[3]: https://pkg.go.dev/github.com/RussellLuo/kok
