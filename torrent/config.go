package torrent

import (
	"github.com/DistributedSolutions/DIMWIT/common/constants"
	"github.com/DistributedSolutions/DIMWIT/util"
	"github.com/anacrolix/torrent"
)

type TopLevelConfig struct {
	AConfig *torrent.Config
}

func NewTopLevelConfig() *TopLevelConfig {
	c := new(TopLevelConfig)
	c.AConfig = &torrent.Config{
		ListenAddr: "0.0.0.0:0",
		DataDir:    util.GetHomeDir() + constants.HIDDEN_DIR + constants.TORRENT_DIR,
	}

	return c
}
