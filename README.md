# kok

The generating toolkit of [Go kit][1].


## Installation

```bash
$ go get -u github.com/RussellLuo/kok/cmd/kok
```

Usage:

```bash
$ kok -h
kok [flags] source-file interface-name
  -fmt
    	whether to make code formatted (default true)
  -out string
    	output directory (default ".")
  -pkg string
    	package name (default will infer)
  -test string
    	the YAML file that provides test-cases for HTTP (default "./http.test.yaml")
```


## Examples

See [examples/profilesvc](examples/profilesvc).


## Documentation

Checkout the [Godoc][2].


## License

[MIT](LICENSE)


[1]: https://github.com/go-kit/kit
[2]: https://godoc.org/github.com/RussellLuo/kok
