package eventcodec

// Codec is a codec (encoder and decoder) for an event.
type Codec interface {
	// Decode decodes an event for subscribers.
	Decode(data, out interface{}) error
}

type Codecs interface {
	EncodeDecoder(name string) Codec
}
