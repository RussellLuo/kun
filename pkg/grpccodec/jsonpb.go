package grpccodec

import (
	"encoding/json"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type ProtoJSON struct{}

func (pj ProtoJSON) DecodeRequest(pb proto.Message, out interface{}) error {
	data, err := protojson.Marshal(pb)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, out)
}

func (pj ProtoJSON) EncodeResponse(in interface{}, pb proto.Message) error {
	data, err := json.Marshal(in)
	if err != nil {
		return err
	}
	return protojson.Unmarshal(data, pb)
}
