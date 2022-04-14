package eventcodec

// NamedCodec holds a codec and its corresponding operation name.
type NamedCodec struct {
	Name  string
	Codec Codec
}

// Op is a shortcut for creating an instance of NamedCodec.
func Op(name string, codec Codec) NamedCodec {
	return NamedCodec{
		Name:  name,
		Codec: codec,
	}
}

type DefaultCodecs struct {
	d      Codec
	Codecs map[string]Codec
}

func NewDefaultCodecs(d Codec, namedCodecs ...NamedCodec) *DefaultCodecs {
	if d == nil {
		d = JSON{} // defaults to JSON
	}

	codecs := make(map[string]Codec)
	for _, c := range namedCodecs {
		codecs[c.Name] = c.Codec
	}

	return &DefaultCodecs{
		d:      d,
		Codecs: codecs,
	}
}

func (dc *DefaultCodecs) EncodeDecoder(name string) Codec {
	if c, ok := dc.Codecs[name]; ok {
		return c
	}
	return dc.d
}
