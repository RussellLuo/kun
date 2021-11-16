package helloworld

import (
	"regexp"

	"github.com/RussellLuo/kok/pkg/httpoption2"
	v "github.com/RussellLuo/validating/v2"
)

var RequestValidators = []httpoption.NamedValidator{
	httpoption.Op("SayHello", ValidateSayHelloRequest(func(req *SayHelloRequest) v.Schema {
		return v.Schema{
			v.F("name", &req.Name): v.All(
				v.Len(0, 10).Msg("length exceeds 10"),
				v.Match(regexp.MustCompile(`^\w+$`)).Msg("invalid name format"),
			),
		}
	})),
}
