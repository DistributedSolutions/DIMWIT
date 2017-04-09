package writeHelper

import (
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
}

func NewWriterHelper(con *constructor.Constructor, fw lite.FactomLiteWriter) *WriteHelper {
	w := new(WriteHelper)
	w.Reader = con
	w.Writer = fw

	return w
}

func (w *WriteHelper) VerifyChannel(ch *common.Channel) (cost int, err *util.ApiError) {
	return 0, util.NewApiError(nil, nil)
}

func (w *WriteHelper) InitiateChannel(ch *common.Channel) (newCh *common.Channel, err *util.ApiError) {
	return
}

func (w *WriteHelper) UpdateChannel(ch *common.Channel) (newCh *common.Channel, err *util.ApiError) {
	return
}

func (w *WriteHelper) DeleteChannel(rootChain *primitives.Hash) (err *util.ApiError) {
	return
}

func (w *WriteHelper) VerifyContent(ch *common.Content) (cost int, err *util.ApiError) {
	return 0, util.NewApiError(nil, nil)
}

func (w *WriteHelper) AddContent(con *common.Content, contentID *primitives.Hash) (chains []*factom.Chain, entries []*factom.Entry, err *util.ApiError) {
	return
}

func (w *WriteHelper) DeleteContent(contentID *primitives.Hash) (err *util.ApiError) {
	return
}
