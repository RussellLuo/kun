package main

import (
	"flag"
	"log"
	"net"

	"github.com/RussellLuo/kun/examples/profilesvcgrpc"
	"github.com/RussellLuo/kun/examples/profilesvcgrpc/pb"
	"github.com/RussellLuo/kun/pkg/grpccodec"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	addr := flag.String("addr", ":8080", "gRPC listen address")
	flag.Parse()

	svc := profilesvcgrpc.NewInmemService()
	server := profilesvcgrpc.NewGRPCServer(svc, grpccodec.NewDefaultCodecs(nil))

	lis, err := net.Listen("tcp", *addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterServiceServer(s, server)
	log.Printf("server listening at %v", lis.Addr())

	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
