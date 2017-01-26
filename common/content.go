package common

import (
	"github.com/DistributedSolutions/DIMWIT/common/primitives"
)

type Content struct {
	// Content Metadada
	Type             byte
	RootChainID      primitives.Hash
	InfoHash         primitives.InfoHash
	LongDescription  primitives.LongDescription
	ShortDescription primitives.ShortDescription
	ActionFiles      primitives.FileList
	Checksum         primitives.MD5Checksum
	Thumbnail        primitives.Image
	Series           byte
	// There can be lots of parts
	Part [2]byte

	// Torrent Metadata
	Tracker  primitives.TrackerList
	FileList primitives.FileList
}
