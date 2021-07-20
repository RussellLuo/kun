package importalias

import (
	dstclient "github.com/RussellLuo/kok/pkg/ifacetool/moq/testdata/importalias/dst/client"
	srcclient "github.com/RussellLuo/kok/pkg/ifacetool/moq/testdata/importalias/src/client"
)

type Connector interface {
	Connect(src srcclient.Client, dst dstclient.Client)
}
