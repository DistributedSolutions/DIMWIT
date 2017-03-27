package torrent

import (
	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/metainfo"
)

type TorrentClient struct {
	client *torrent.Client
}

func NewTorrentClientFromConfig(con *TopLevelConfig) *TorrentClient {
	c := new(TorrentClient)
	c.client = torrent.NewClient(con.AConfig)

	return c
}

func NewTorrentClient() *TorrentClient {
	return NewTorrentClientFromConfig(NewTopLevelConfig())
}

func (c *TorrentClient) AddMagnet(mag *metainfo.Magnet) (*torrent.Torrent, error) {
	return c.client.AddMagnet(mag.String())
}
