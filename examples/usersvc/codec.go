package usersvc

import (
	"fmt"
	"net"
	"strings"

	"github.com/RussellLuo/kun/pkg/httpcodec"
)

// IPCodec is used to encode and decode an IP. It can be reused wherever needed.
type IPCodec struct{}

func (c IPCodec) Decode(in map[string][]string, out interface{}) error {
	remote := in["request.RemoteAddr"][0]

	fwdFor := in["header.X-Forwarded-For"]
	if len(fwdFor) > 0 && fwdFor[0] != "" {
		// Prefer the first IP.
		remote = strings.TrimSpace(strings.Split(fwdFor[0], ",")[0])
	}

	ipStr, _, err := net.SplitHostPort(remote)
	if err != nil {
		ipStr = remote // OK; probably didn't have a port
	}

	ip := net.ParseIP(ipStr)
	if ip == nil {
		return fmt.Errorf("invalid client IP address: %s", ipStr)
	}

	outIP := out.(*net.IP)
	*outIP = ip
	return nil
}

func (c IPCodec) Encode(in interface{}) (out map[string][]string) {
	return nil
}

func NewCodecs() *httpcodec.DefaultCodecs {
	// Use IPCodec to encode and decode the "IP" field in the struct argument
	// named "user", if exists, for the operation named "CreateUser".
	return httpcodec.NewDefaultCodecs(nil,
		httpcodec.Op("CreateUser", httpcodec.NewPatcher(httpcodec.JSON{}).Params(
			"user", httpcodec.StructParams{
				Fields: map[string]httpcodec.ParamsCodec{
					"IP": IPCodec{},
				},
			},
		)))
}
