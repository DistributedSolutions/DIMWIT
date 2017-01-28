package common

import (
	"github.com/DistributedSolutions/DIMWIT/common/primitives"
)

type Channel struct {
	RootChainID       primitives.Hash
	ManagementChainID primitives.Hash
	ContentChainID    primitives.Hash
	// They are not an array, because they are never referenced as an array
	LV1PublicKey      primitives.PublicKey
	LV2PublicKey      primitives.PublicKey
	LV3PublicKey      primitives.PublicKey
	ContentSingingKey primitives.PublicKey

	ChannelTitle     primitives.Title
	Website          primitives.SiteURL
	LongDescription  primitives.LongDescription
	ShortDescription primitives.ShortDescription
	Playlist         ManyPlayList
	Thumbnail        primitives.Image
	Banner           primitives.Image
	Tags             primitives.TagList
	SuggestedChannel primitives.HashList
	Content          ContentList
}
