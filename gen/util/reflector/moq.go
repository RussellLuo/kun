package reflector

import (
	"github.com/RussellLuo/kok/pkg/ifacetool"
	"github.com/RussellLuo/kok/pkg/ifacetool/moq"
)

func ParseInterface(srcDir, pkgName, ifaceName string) (*ifacetool.Data, error) {
	parser, err := moq.New(moq.Config{SrcDir: srcDir, PkgName: pkgName})
	if err != nil {
		return nil, err
	}

	data, err := parser.Parse(ifaceName)
	if err != nil {
		return nil, err
	}

	return data, nil
}
