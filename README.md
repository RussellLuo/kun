# kok

The toolkit of [Go kit][1].


## Zen

- Just write code for your business logic, generate everything else.
- Implement the service once, consume it in various ways (in-process function call or RPC).


## Features

1. Code Generation Tool

    - [x] Endpoint
    - [x] Transport
        - [x] HTTP
            + [x] HTTP Server
            + [x] HTTP Test
            + [x] HTTP Client
            + [x] [OAS-v2][2] Documentation
        - [ ] gRPC

2. Useful Packages

    - [appx](pkg/appx): Application framework for HTTP and CRON applications (a wrapper of [appx][3]).
    - [prometheus](pkg/prometheus): Prometheus metrics utilities.
    - [trace](pkg/trace): A thin wrapper of [x/net/trace][4] for Go kit.
    - [werror](pkg/werror): Classified business errors.


## Installation

```bash
$ go get -u github.com/RussellLuo/kok/cmd/kokgen
```

<details open>
  <summary> Usage </summary>

```bash
$ kokgen -h
kokgen [flags] source-file interface-name
  -flat
        whether to use flat layout (default true)
  -fmt
        whether to make code formatted (default true)
  -out string
        output directory (default ".")
  -snake
        whether to use snake-case for default names (default true)
  -test string
        the YAML file that provides test-cases for HTTP (default "./http.test.yaml")
  -trace
        whether to enable tracing
```

</details>


## Quick Start

**NOTE**: The following code is located in [helloworld](examples/helloworld).

1. Define the interface

    ```go
    type Service interface {
        SayHello(ctx context.Context, name string) (message string, err error)
    }
    ```

2. Implement the service

    ```go
    type Greeter struct{}

    func (g *Greeter) SayHello(ctx context.Context, name string) (string, error) {
        return "Hello " + name, nil
    }
    ```

3. Add HTTP annotations

    ```go
    type Service interface {
        // @kok(op): POST /messages
        SayHello(ctx context.Context, name string) (message string, err error)
    }
    ```

4. Generate the HTTP code

    ```bash
    $ cd examples/helloworld
    $ kokgen ./service.go Service
    ```

5. Consume the service

    Run the HTTP server:

    ```bash
    $ go run cmd/main.go
    2020/09/15 18:06:22 transport=HTTP addr=:8080
    ```

    Consume by [HTTPie](https://github.com/jakubroztocil/httpie):

    ```bash
    $ http POST :8080/messages name=Tracey
    HTTP/1.1 200 OK
    Content-Length: 27
    Content-Type: application/json; charset=utf-8
    Date: Tue, 15 Sep 2020 10:06:34 GMT

    {
        "message": "Hello Tracey"
    }
    ```

6. See the OAS documentation

    <details>
      <summary> (Click to expand) </summary>

    ```bash
    $ http GET :8080/api
    HTTP/1.1 200 OK
    Content-Length: 848
    Content-Type: text/plain; charset=utf-8
    Date: Tue, 15 Sep 2020 10:08:24 GMT

    swagger: "2.0"
    info:
      title: "No Title"
      version: "0.0.0"
      description: "Service is used for saying hello."
      license:
        name: "MIT"
    host: "example.com"
    basePath: "/"
    schemes:
      - "https"
    consumes:
      - "application/json"
    produces:
      - "application/json"

    paths:
      /messages:
        post:
          description: "SayHello says hello to the given name."
          operationId: "SayHello"
          parameters:
            - name: body
              in: body
              schema:
                $ref: "#/definitions/SayHelloRequestBody"

          produces:
            - application/json; charset=utf-8
          responses:
            200:
              description: ""
              schema:
                $ref: "#/definitions/SayHelloResponse"


    definitions:
      SayHelloRequestBody:
        type: object
        properties:
          name:
            type: string
      SayHelloResponse:
        type: object
        properties:
          message:
            type: string
    ```

    </details>

See more examples [here](examples).


## HTTP

### Annotations

<details open>
  <summary> Define the HTTP request operation </summary>

- Key: `@kok(op)`
- Value: `<method> <pattern>`
    + **method**: The request method
    + **pattern**: The request URL
        - NOTE: All variables in **pattern** will automatically be bound to their corresponding method arguments (matches by name in *lower camel case*), as **path** parameters, if the variables are not yet specified as path parameters explicitly by `@kok(param)`.
- Example:

    ```go
    type Service interface {
        // @kok(op): DELETE /users/{id}
        DeleteUser(ctx context.Context, id int) (err error)
    }

    // HTTP request:
    // $ http DELETE /users/101
    ```

</details>

<details open>
  <summary> Define the HTTP request parameters </summary>

- Key: `@kok(param)`
- Value: `<argName> < in:<in>,name:<name>,required:<required>`
    + **argName**: The name of the method argument.
        - *Argument aggregation*: By specifying the same **argName**, multiple request parameters (each one is of basic type or repeated basic type) can be aggregated into one method argument (of any type).
            + You do not need to repeat the **argName**, only the first one is required.
    + **in**:
        - **path**: The method argument is sourced from a [path parameter](https://swagger.io/docs/specification/describing-parameters/#path-parameters).
            + Optional: All variables in **pattern** will automatically be bound to their corresponding method arguments (matches by name in *lower camel case*), as **path** parameters.
        - **query**: The method argument is sourced from a [query parameter](https://swagger.io/docs/specification/describing-parameters/#query-parameters).
            + To receive values from a multi-valued query parameter, the method argument can be defined as a slice of basic type.
        - **header**: The method argument is sourced from a [header parameter](https://swagger.io/docs/specification/describing-parameters/#header-parameters).
        - **cookie**: The method argument is sourced from a [cookie parameter](https://swagger.io/docs/specification/describing-parameters/#cookie-parameters).
            + Not supported yet.
        - **request**: The method argument is sourced from a property of Go's [http.Request](https://golang.org/pkg/net/http/#Request).
            + This is a special case, and only one property `RemoteAddr` is available now.
            + Note that parameters located in **request** have no relationship with OAS.
    + **name**: The name of the corresponding request parameter.
        - Optional: Defaults to **argName** (snake-case, or lower-camel-case if `-snake=false`) if not specified.
    + **descr**: The OAS description of the corresponding request parameter.
        - Optional: Defaults to `""` if not specified.
    + **required**: Determines whether this parameter is mandatory.
        - Optional: Defaults to `false`, if not specified.
        - If the parameter location is **path**, this property will be set to `true` internally, whether it's specified or not.
- Example:
    + Bind request parameters to simple arguments:

        ```go
        type Service interface {
            // @kok(op): PUT /users/{id}
            // @kok(param): name < in:header,name:X-User-Name
            UpdateUser(ctx context.Context, id int, name string) (err error)
        }

        // HTTP request:
        // $ http PUT /users/101 X-User-Name:tracey
        ```
    + Bind multiple request parameters to a struct according to tags:

        ```go
        type User struct {
            ID   int    `kok:"path.id"`
            Name string `kok:"query.name"`
            Age  int    `kok:"header.X-User-Age"`
        }

        type Service interface {
            // @kok(op): PUT /users/{id}
            // @kok(param): user
            UpdateUser(ctx context.Context, user User) (err error)
        }

        // HTTP request:
        // $ http PUT /users/101?name=tracey X-User-Age:1
        ```
    + Bind multiple query parameters to a struct with no tags:

        ```go
        type User struct {
            Name    string
            Age     int
            Hobbies []string
        }

        type Service interface {
            // @kok(op): POST /users
            // @kok(param): user
            CreateUser(ctx context.Context, user User) (err error)
        }

        // HTTP request:
        // $ http POST /users?Name=tracey&Age=1&Hobbies=music&Hobbies=sport
        ```
    + Argument aggregation:

        ```go
        type Service interface {
            // @kok(op): POST /logs
            // @kok(param): ip < in:header,name:X-Forwarded-For
            // @kok(param): ip < in:request,name:RemoteAddr
            Log(ctx context.Context, ip net.IP) (err error)
        }

        // The equivalent annotations.
        type Service interface {
            // @kok(op): POST /logs
            // @kok(param): ip < in:header,name:X-Forwarded-For
            // @kok(param):    < in:request,name:RemoteAddr
            Log(ctx context.Context, ip net.IP) (err error)
        }

        // You must customize the decoding of `ip` later (conventionally in another file named `codec.go`).
        // See examples in the `Encoding and decoding` section.

        // HTTP request:
        // $ http POST /logs
        ```

</details>

<details open>
  <summary> Define the HTTP request body </summary>

- Key: `@kok(body)`
- Value: `<field>` or `body:<field>,name:<argName>=<name>,descr:<argName>=<descr>`
    + **field**: The name of the method argument whose value is mapped to the HTTP request body.
        - Optional: When omitted, a struct containing all the arguments, which are not located in **path**/**query**/**header**, will automatically be mapped to the HTTP request body.
        - The special name `-` can be used, to define that there is no HTTP request body. As a result, every argument, which is not located in **path**/**query**/**header**, will automatically be mapped to one or more query parameters.
    + **argName**: The name of the method argument to be manipulated.
    + **name**: The name of the corresponding request parameter.
        - Optional: Defaults to **argName** (snake-case, or lower-camel-case if `-snake=false`) if not specified.
    + **descr**: The OAS description of the corresponding request parameter.
        - Optional: Defaults to `""` if not specified.
- Example:
    + Omitted:

        ```go
        type Service interface {
            // @kok(op): POST /users
            CreateUser(ctx context.Context, name string, age int) (err error)
        }

        // HTTP request:
        // $ http POST /users name=tracey age=1
        ```

    + Specified as a normal argument:

        ```go
        type User struct {
            Name string `json:"name"`
            Age  int    `json:"age"`
        }

        type Service interface {
            // @kok(op): POST /users
            // @kok(body): user
            CreateUser(ctx context.Context, user User) (err error)
        }

        // HTTP request:
        // $ http POST /users name=tracey age=1
        ```

    + Specified as `-`:

        ```go
        type User struct {
            Name    string   `kok:"name"`
            Age     int      `kok:"age"`
            Hobbies []string `kok:"hobby"`
        }

        type Service interface {
            // @kok(op): POST /users
            // @kok(body): -
            CreateUser(ctx context.Context, user User) (err error)
        }

        // HTTP request:
        // $ http POST /users?name=tracey&age=1&hobby=music&hobby=sport
        ```

</details>

<details open>
  <summary> Define the success HTTP response </summary>


- Key: `@kok(success)`
- Value: `statusCode:<statusCode>,body:<body>`
    + **statusCode**: The status code of the success HTTP response.
        - Optional: Defaults to 200 if not specified.
    + **body**: The name of the response field whose value is mapped to the HTTP response body.
        - Optional: When omitted, a struct containing all the results (except error) will automatically be mapped to the HTTP response body.
- Example:

    ```go
    type User struct {
        Name string `json:"name"`
        Age  int    `json:"age"`
    }

    type Service interface {
        // @kok(op): POST /users
        // @kok(success): statusCode:201,body:user
        CreateUser(ctx context.Context) (user User, err error)
    }
    ```

</details>

### Encoding and decoding

See the [HTTP Codec](https://github.com/RussellLuo/kok/blob/master/pkg/codec/httpcodec/codec.go#L8-L22) interface.

Also see [here](https://github.com/RussellLuo/kok/issues/8) for examples.


### OAS Schema

See the [OAS Schema](https://github.com/RussellLuo/kok/blob/master/pkg/oasv2/schema.go#L18-L21) interface.


## Documentation

Checkout the [Godoc][5].


## License

[MIT](LICENSE)


[1]: https://github.com/go-kit/kit
[2]: https://swagger.io/specification/v2/
[3]: https://github.com/RussellLuo/appx
[4]: https://pkg.go.dev/golang.org/x/net/trace
[5]: https://pkg.go.dev/github.com/RussellLuo/kok
