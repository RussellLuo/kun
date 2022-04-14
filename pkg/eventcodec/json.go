package eventcodec

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/RussellLuo/kun/pkg/eventpubsub"
	"github.com/RussellLuo/kun/pkg/werror"
)

type JSON struct{}

func (j JSON) Decode(data, out interface{}) error {
	switch d := data.(type) {
	case []byte:
		if err := json.Unmarshal(d, out); err != nil {
			return werror.Wrap(eventpubsub.ErrInvalidData, err)
		}
	case io.Reader:
		if err := json.NewDecoder(d).Decode(out); err != nil {
			return werror.Wrap(eventpubsub.ErrInvalidData, err)
		}
	default:
		return werror.Wrap(eventpubsub.ErrInvalidData, fmt.Errorf("unsupported data type %T", d))
	}
	return nil
}
