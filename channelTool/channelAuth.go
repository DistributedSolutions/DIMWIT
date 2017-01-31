package channelTool

import (
	"github.com/DistributedSolutions/DIMWIT/common"
	"github.com/DistributedSolutions/DIMWIT/common/primitives"
)

type AuthChannel struct {
	Channel common.Channel

	PrivateKeys [3]primitives.PrivateKey
}
