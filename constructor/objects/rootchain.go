package objects

import (
	//"github.com/FactomProject/factom"
	//"github.com/DistributedSolutions/DIMWIT/common"
	"github.com/DistributedSolutions/DIMWIT/common/primitives"
	"github.com/DistributedSolutions/DIMWIT/factom-lite"
)

// Factom Chain
//		byte		Version
//		[18]byte	"Channel Root Chain"
//		[]byte		Title
//		[32]byte	PublicKey(1)
//		[32]byte	PublicKey(2)
//		[32]byte	PublicKey(3)
//		[]byte		Nonce
type RootChainApplyEntry struct {
	// Memory
	Channel *ChannelWrapper
	Entry   *lite.EntryHolder

	//  Object
	Version byte
	Title   primitives.Title
	PubKeys []primitives.PublicKey
}

func NewRootChainApplyEntry() IApplyEntry {
	m := new(RootChainApplyEntry)
	m.PubKeys = make([]primitives.PublicKey, 3)
	return m
}

func (r *RootChainApplyEntry) ParseFactomEntry(e *lite.EntryHolder) error {
	ent := e.Entry
	r.Version = ent.ExtIDs[0][0]
	tit := new(primitives.Title)
	err := tit.UnmarshalBinary(ent.ExtIDs[2])
	if err != nil {
		return err
	}
	r.Title = *tit

	for i := 0; i < 3; i++ {
		p, err := primitives.PublicKeyFromBytes(ent.ExtIDs[3+i])
		if err != nil {
			return err
		}
		r.PubKeys[i] = *p
	}

	r.Entry = e
	return nil
}

func (r *RootChainApplyEntry) RequestChannel() (string, bool) {
	return r.Entry.Entry.ChainID, false
}

func (r *RootChainApplyEntry) AnswerChannelRequest(cw *ChannelWrapper) error {
	r.Channel = cw
	return nil
}

func (r *RootChainApplyEntry) NeedIsFirstEntry() bool { return true }

func (r *RootChainApplyEntry) ApplyEntry() (*ChannelWrapper, bool) {
	cw := new(ChannelWrapper)
	chainID, err := primitives.HexToHash(r.Entry.Entry.ChainID)
	if err != nil {
		return nil, false
	}

	// Instantiate Channel
	cw.Channel.RootChainID = *chainID
	cw.RMadeHeight = r.Entry.Height
	cw.Channel.ChannelTitle = r.Title
	cw.Channel.LV1PublicKey = r.PubKeys[0]
	cw.Channel.LV2PublicKey = r.PubKeys[1]
	cw.Channel.LV3PublicKey = r.PubKeys[2]

	return cw, false
}

// Unused
func (m *RootChainApplyEntry) RequestEntriesInOtherChain() (string, bool)         { return "", false }
func (m *RootChainApplyEntry) AnswerChainEntriesInOther(ents []*lite.EntryHolder) {}
func (r *RootChainApplyEntry) NeedChainEntries() bool                             { return false }
func (r *RootChainApplyEntry) AnswerChainEntries(ents []*lite.EntryHolder)        {}
