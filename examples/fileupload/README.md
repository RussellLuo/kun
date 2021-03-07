# fileupload

This example illustrates how to handle file uploading.


## Generate the code

```bash
$ go generate
```

## Test the server

Run the server:

```bash
$ go run cmd/main.go
2021/03/07 16:00:34 transport=HTTP addr=:8080
```

Upload a file:

```bash
$ curl -F file=@/dir/to/sample_file.txt http://localhost:8080/upload
```

Check the uploaded file `sample_file.txt` in the current directory.
