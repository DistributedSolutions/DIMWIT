package lite

import (
	"github.com/FactomProject/factom"
)

// The calls that can be made to the factomlite client
type FactomLite interface {
	// Does the commit + reveal
	SubmitEntry(e factom.Entry, ec factom.ECAddress) error
	SubmitChain(e factom.Chain, ec factom.ECAddress) error
}
