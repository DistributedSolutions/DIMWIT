package elements

type IElement interface {
}

type IElementPiece interface {
	// Ascii Type in bytes
	Type() []byte
	IsChain() bool
	ForChain() int
}

// Types
var (
	// Root
	TYPE_ROOT_CHAIN       = []byte("Channel Root Chain")
	TYPE_ROOT_REGISTER    = []byte("Channel Chain")
	TYPE_ROOT_CONTENT_KEY = []byte("Content Signing Key")

	// Manage
	TYPE_MANAGE_CHAIN          = []byte("Channel Management Chain")
	TYPE_MANAGE_CHAIN_REGISTER = []byte("Register Management Chain")
	TYPE_MANAGE_CHAIN_METADATA = []byte("Channel Management Metadata Main")

	// ChanContentChain
	TYPE_CHANNEL_CONENT_CHAIN          = []byte("Channel Content Chain")
	TYPE_CHANNEL_CONENT_CHAIN_REGISTER = []byte("Register Content Chain")

	// SingleContent
	TYPE_SINGLE_CONTENT_CHAIN = []byte("Content Chain")
	TYPE_CONTENT_LINK         = []byte("Content Link")
)

// Chains
const (
	CHAIN_NA int = iota
	CHAIN_MAIN
	CHAIN_ROOT
	CHAIN_MANAGEMENT
	CHAIN_CONTENT_LIST
	CHAIN_CONTENT_SINGLE
)
