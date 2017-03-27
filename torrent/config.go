package torrent

import (
	"github.com/anacrolix/torrent"
)

type TopLevelConfig struct {
	AConfig *torrent.Config
}

func NewTopLevelConfig() *TopLevelConfig {
	c := new(TopLevelConfig)
	c.AConfig = torrent.Config{}

	return c
}
