package engine

import (
	"github.com/DistributedSolutions/DIMWIT/constructor"
	"github.com/DistributedSolutions/DIMWIT/factom-lite"
	"github.com/DistributedSolutions/DIMWIT/provider"
	"github.com/DistributedSolutions/DIMWIT/torrent"
)

type WholeState struct {
	Constructor   *constructor.Constructor
	Provider      *provider.Provider
	FactomClient  lite.FactomLite
	TorrentClient *torrent.TorrentClient
}
