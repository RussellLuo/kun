package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/RussellLuo/kok/examples/fileupload"
	"github.com/RussellLuo/kok/pkg/codec/httpcodec"
	"github.com/RussellLuo/kok/pkg/oasv2"
)

func main() {
	httpAddr := flag.String("http.addr", ":8080", "HTTP listen address")
	flag.Parse()

	svc := &fileupload.Uploader{}
	codecs := httpcodec.CodecMap{Default: httpcodec.NewMultipartForm(0)}
	r := fileupload.NewHTTPRouterWithOAS(svc, codecs, &oasv2.ResponseSchema{})

	errs := make(chan error, 2)
	go func() {
		log.Printf("transport=HTTP addr=%s\n", *httpAddr)
		errs <- http.ListenAndServe(*httpAddr, r)
	}()
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	log.Printf("terminated, err:%v", <-errs)
}
