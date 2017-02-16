package lite

import (
	"github.com/DistributedSolutions/DIMWIT/common/primitives"
	"github.com/FactomProject/factom"
)

// The calls that can be made to the factomlite client
type FactomLite interface {
	FactomLiteReader
	FactomLiteWriter
}

// FactomLiteReader is the readonly functions of a FactomLite Client
type FactomLiteReader interface {
	GetAllChainEntries(chainID primitives.Hash) ([]*factom.Entry, error)
	GetFirstEntry(chainID primitives.Hash) (*factom.Entry, error)
	GetEntry(entryHash primitives.Hash) (*factom.Entry, error)
	GetReadyHeight() (uint32, error)
	GrabAllEntriesAtHeight(height uint32) ([]*EntryHolder, error)
}

// FactomLiteWriter is the writeonly functions of a FactomLite Client
type FactomLiteWriter interface {
	// Does the commit + reveal
	SubmitEntry(e factom.Entry, ec factom.ECAddress) error
	SubmitChain(e factom.Chain, ec factom.ECAddress) error
}
