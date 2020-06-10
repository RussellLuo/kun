# kok

The toolkit of [Go kit][1].


## Installation

```bash
$ go get -u github.com/RussellLuo/kok/cmd/kokgen
```

Usage:

```bash
$ kokgen -h
kok [flags] source-file interface-name
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


## Examples

See [examples/profilesvc](examples/profilesvc).


## Documentation

Checkout the [Godoc][2].


## License

[MIT](LICENSE)


[1]: https://github.com/go-kit/kit
[2]: https://godoc.org/github.com/RussellLuo/kok
