package lite

import (
	"encoding/hex"
	"fmt"
	"sync"
	"time"

	//"github.com/DistributedSolutions/DIMWIT/common/constants"
	"github.com/DistributedSolutions/DIMWIT/common/primitives"
	"github.com/DistributedSolutions/DIMWIT/database"
	"github.com/FactomProject/factom"
)

var (
	FACTOM_BUCKET []byte = []byte("FactomShit")
)

// A very dumb factom-lite client
// It acts as both a Writer, and a Reader
type FakeDumbLite struct {
	FactomdLocation string

	height     uint32
	heightlist [][]factom.Entry // Entries at given height
	chainlists map[string][]factom.Entry

	db database.IDatabase
	sync.RWMutex
}

func NewFakeDumbLite() FactomLite {
	d := new(FakeDumbLite)
	d.db = database.NewMapDB()
	d.chainlists = make(map[string][]factom.Entry)
	d.heightlist = make([][]factom.Entry, 50000)
	return d
}

//
// Writer Interface Methods
//

func (d *FakeDumbLite) SubmitEntry(e factom.Entry, ec factom.ECAddress) (comId string, eHash string, err error) {
	data, err := e.MarshalJSON()
	if err != nil {
		return "", "", err
	}
	d.db.Put(FACTOM_BUCKET, e.Hash(), data)

	d.Lock()
	d.heightlist[d.height] = append(d.heightlist[d.height], e)
	d.chainlists[e.ChainID] = append(d.chainlists[e.ChainID], e)
	d.Unlock()
	d.height++
	return "", hex.EncodeToString(e.Hash()), nil
}

func (d *FakeDumbLite) SubmitChain(c factom.Chain, ec factom.ECAddress) (comId string, chainID string, err error) {
	e := c.FirstEntry
	data, err := e.MarshalJSON()
	if err != nil {
		return "", "", err
	}
	d.db.Put(FACTOM_BUCKET, e.Hash(), data)

	d.Lock()
	d.heightlist[d.height] = append(d.heightlist[d.height], *e)
	d.chainlists[e.ChainID] = append([]factom.Entry{*e}, d.chainlists[e.ChainID]...)
	d.Unlock()
	d.height++
	return "", c.FirstEntry.ChainID, nil
}

//
// Reader Interface Methods
//

func (d *FakeDumbLite) GetAllChainEntries(chainID primitives.Hash) ([]*factom.Entry, error) {
	d.RLock()
	defer d.RUnlock()
	ents := d.chainlists[chainID.String()]
	entries := make([]*factom.Entry, 0)
	for _, e := range ents {
		data, err := d.db.Get(FACTOM_BUCKET, e.Hash())
		if err != nil {
			return nil, err
		}
		en := new(factom.Entry)
		err = en.UnmarshalJSON(data)
		if err != nil {
			return nil, err
		}
		entries = append(entries, en)
	}
	return entries, nil
}

func (d *FakeDumbLite) GetFirstEntry(chainID primitives.Hash) (*factom.Entry, error) {
	d.RLock()
	defer d.RUnlock()
	ents := d.chainlists[chainID.String()]
	e := ents[0]
	data, err := d.db.Get(FACTOM_BUCKET, e.Hash())
	if err != nil {
		return nil, err
	}
	en := new(factom.Entry)
	err = en.UnmarshalJSON(data)
	if err != nil {
		return nil, err
	}
	return en, nil
}

func (d *FakeDumbLite) GetEntry(entryHash primitives.Hash) (*factom.Entry, error) {
	d.RLock()
	defer d.RUnlock()
	data, err := d.db.Get(FACTOM_BUCKET, entryHash.Bytes())
	if err != nil {
		return nil, err
	}
	en := new(factom.Entry)
	err = en.UnmarshalJSON(data)
	if err != nil {
		return nil, err
	}
	return en, nil
}

func (d *FakeDumbLite) GetReadyHeight() (uint32, error) {
	d.RLock()
	defer d.RUnlock()
	return d.height, nil
}

func (d *FakeDumbLite) GrabAllEntriesAtHeight(height uint32) ([]*EntryHolder, error) {
	d.RLock()
	defer d.RUnlock()
	rh, err := d.GetReadyHeight()
	if err != nil {
		return nil, err
	}
	if height > rh {
		return nil, fmt.Errorf("The height given is not ready to be grabbed and parsed. Given %d, ready up to %d", height, rh)
	}

	entries := make([]*EntryHolder, 0)
	// Cycle through entries
	ents := d.heightlist[height]
	for _, e := range ents {
		eholder := new(EntryHolder)
		eholder.Timestamp = time.Now().Unix()
		eholder.Entry = &e
		eholder.Height = height

		entries = append(entries, eholder)
	}

	return entries, nil
}
