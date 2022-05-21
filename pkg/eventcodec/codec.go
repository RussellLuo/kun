package eventcodec

// Codec is a codec (encoder and decoder) for an event.
type Codec interface {
	// Decode decodes data and stores the result in out.
	Decode(data, out interface{}) error

	// Encode encodes in and stores the result in data.
	Encode(in interface{}) (data interface{}, err error)
}

type Codecs interface {
	EncodeDecoder(name string) Codec
}
