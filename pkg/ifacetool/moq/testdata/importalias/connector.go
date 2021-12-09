package importalias

import (
	dstclient "github.com/RussellLuo/kun/pkg/ifacetool/moq/testdata/importalias/dst/client"
	srcclient "github.com/RussellLuo/kun/pkg/ifacetool/moq/testdata/importalias/src/client"
)

type Connector interface {
	Connect(src srcclient.Client, dst dstclient.Client)
}
