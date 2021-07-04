package grpccodec

import (
	"google.golang.org/protobuf/proto"
)

// Codec is a series of codecs (encoders and decoders) for gRPC requests and responses.
type Codec interface {
	// DecodeRequest converts a proto message to a Go value.
	// It is designed to be used at the server side.
	DecodeRequest(pb proto.Message, out interface{}) error

	// EncodeResponse converts a Go value to a proto message.
	// It is designed to be used at the server side.
	EncodeResponse(in interface{}, pb proto.Message) error
}

type Codecs interface {
	EncodeDecoder(name string) Codec
}
