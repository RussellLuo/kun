package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/RussellLuo/kun/examples/profilesvc"
)

func main() {
	baseURL := flag.String("url", "http://localhost:8080", "The base URL")
	flag.Parse()

	client, err := profilesvc.NewHTTPClient(
		profilesvc.NewCodecs(),
		&http.Client{Timeout: 10 * time.Second},
		*baseURL,
	)
	if err != nil {
		log.Fatalf("NewHTTPClient err: %v\n", err)
	}

	if err = client.PostProfile(context.Background(), profilesvc.Profile{ID: "1", Name: "profile1"}); err != nil {
		log.Fatalf("PostProfile err: %v\n", err)
	}

	profile, err := client.GetProfile(context.Background(), "1")
	if err != nil {
		log.Fatalf("GetProfile err: %v\n", err)
	}
	log.Printf("GetProfile ok: %+v\n", profile)

	if err := client.PutProfile(context.Background(), "1", profilesvc.Profile{ID: "1", Name: "profile2"}); err != nil {
		log.Fatalf("PutProfile err: %v\n", err)
	}

	if err := client.PatchProfile(context.Background(), "1", profilesvc.Profile{ID: "1", Name: "profile3"}); err != nil {
		log.Fatalf("PatchProfile err: %v\n", err)
	}

	if err = client.PostAddress(context.Background(), "1", profilesvc.Address{ID: "4", Location: "address4"}); err != nil {
		log.Fatalf("PostAddress err: %v\n", err)
	}

	address, err := client.GetAddress(context.Background(), "1", "4")
	if err != nil {
		log.Fatalf("GetAddress err: %v\n", err)
	}
	log.Printf("GetAddress ok: %+v\n", address)

	addresses, err := client.GetAddresses(context.Background(), "1")
	if err != nil {
		log.Fatalf("GetAddresses err: %v\n", err)
	}
	log.Printf("GetAddresses ok: %+v", addresses)

	if err := client.DeleteAddress(context.Background(), "1", "4"); err != nil {
		log.Fatalf("DeleteAddress err: %v\n", err)
	}

	if err := client.DeleteProfile(context.Background(), "1"); err != nil {
		log.Fatalf("DeleteProfile err: %v\n", err)
	}
}
