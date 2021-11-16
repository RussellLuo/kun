package httpcodec

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

// PatchAll patches the default codec and all the custom codecs.
func (dc *DefaultCodecs) PatchAll(patch func(Codec) *Patcher) *DefaultCodecs {
	if dc.d != nil {
		dc.d = patch(dc.d)
	}

	for name, c := range dc.Codecs {
		dc.Codecs[name] = patch(c)
	}

	return dc
}

func (dc *DefaultCodecs) EncodeDecoder(name string) Codec {
	if c, ok := dc.Codecs[name]; ok {
		return c
	}
	return dc.d
}

// Patcher is used to change the encoding and decoding behaviors of an
// existing instance of Codec.
type Patcher struct {
	Codec // the original Codec

	// Custom codecs each for a single request parameter.
	paramCodecs map[string]ParamCodec
	// Custom codecs each for a group of request parameters.
	paramsCodecs map[string]ParamsCodec
}

func NewPatcher(codec Codec) *Patcher {
	return &Patcher{
		Codec:        codec,
		paramCodecs:  make(map[string]ParamCodec),
		paramsCodecs: make(map[string]ParamsCodec),
	}
}

// Param sets a codec for a request parameter specified by name.
func (p *Patcher) Param(name string, pc ParamCodec) *Patcher {
	p.paramCodecs[name] = pc
	return p
}

// Params sets a codec for a group of request parameters specified by name.
func (p *Patcher) Params(name string, psc ParamsCodec) *Patcher {
	p.paramsCodecs[name] = psc
	return p
}

func (p *Patcher) DecodeRequestParam(name string, values []string, out interface{}) error {
	if c, ok := p.paramCodecs[name]; ok {
		return c.Decode(values, out)
	}
	return p.Codec.DecodeRequestParam(name, values, out)
}

func (p *Patcher) DecodeRequestParams(name string, values map[string][]string, out interface{}) error {
	if c, ok := p.paramsCodecs[name]; ok {
		return c.Decode(values, out)
	}
	return p.Codec.DecodeRequestParams(name, values, out)
}

func (p *Patcher) EncodeRequestParam(name string, value interface{}) []string {
	if c, ok := p.paramCodecs[name]; ok {
		return c.Encode(value)
	}
	return p.Codec.EncodeRequestParam(name, value)
}

func (p *Patcher) EncodeRequestParams(name string, value interface{}) map[string][]string {
	if c, ok := p.paramsCodecs[name]; ok {
		return c.Encode(value)
	}
	return p.Codec.EncodeRequestParams(name, value)
}
