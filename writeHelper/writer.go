package writeHelper

import (
	"fmt"

	"github.com/DistributedSolutions/DIMWIT/common"
	"github.com/DistributedSolutions/DIMWIT/common/primitives"
	"github.com/DistributedSolutions/DIMWIT/constructor"
	"github.com/DistributedSolutions/DIMWIT/factom-lite"
	"github.com/DistributedSolutions/DIMWIT/util"
	"github.com/FactomProject/factom"
)

type WriteHelper struct {
	// To write into Factom
	Writer lite.FactomLiteWriter

	// To read from Factom
	Reader *constructor.Constructor

	// Map of AuthChannels
	AuthChannels map[string]*AuthChannel

	// ECKey
	ECAddress *factom.ECAddress
}

func NewWriterHelper(con *constructor.Constructor, fw lite.FactomLiteWriter) (*WriteHelper, error) {
	w := new(WriteHelper)
	w.Reader = con
	w.Writer = fw

	w.AuthChannels = make(map[string]*AuthChannel)

	pk, err := primitives.GeneratePrivateKey()
	if err != nil {
		return nil, err
	}

	w.ECAddress, err = factom.MakeECAddress(pk.Secret[:32])
	if err != nil {
		return nil, err
	}

	return w, nil
}

func (w *WriteHelper) VerifyChannel(ch *common.Channel) (cost int, apiErr *util.ApiError) {
	return 0, util.NewAPIError(nil, nil)
}

func (w *WriteHelper) InitiateChannel(ch *common.Channel) (apiErr *util.ApiError) {
	// TODO: Check Balance

	// Generate Keys
	a, err := NewAuthChannel(ch, w.ECAddress)
	if err != nil {
		return util.NewAPIError(err, fmt.Errorf("failed to generate channel keys"))
	}

	// Brute force a ChainID
	a.Initiate(ch)

	// Get Factom Elements
	root := a.Root
	entries, chain := root.FactomElements()

	//	Enter into Factom
	w.Writer.SubmitChain(*chain, *w.ECAddress)
	for _, e := range entries {
		w.Writer.SubmitEntry(*e, *w.ECAddress)
	}

	// Add to our Map
	w.AuthChannels[ch.RootChainID.String()] = a
	return nil
}

func (w *WriteHelper) UpdateChannel(ch *common.Channel) (newCh *common.Channel, apiErr *util.ApiError) {
	return
}

func (w *WriteHelper) DeleteChannel(rootChain *primitives.Hash) (apiErr *util.ApiError) {
	return
}

func (w *WriteHelper) VerifyContent(ch *common.Content) (cost int, apiErr *util.ApiError) {
	return 0, util.NewAPIError(nil, nil)
}

func (w *WriteHelper) AddContent(con *common.Content, contentID *primitives.Hash) (chains []*factom.Chain, entries []*factom.Entry, apiErr *util.ApiError) {
	return
}

func (w *WriteHelper) DeleteContent(contentID *primitives.Hash) (apiErr *util.ApiError) {
	return
}
