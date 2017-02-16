package lite

import (
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/DistributedSolutions/DIMWIT/common/constants"
	"github.com/DistributedSolutions/DIMWIT/common/primitives"
	"github.com/FactomProject/factom"
)

// A very dumb factom-lite client
// It acts as both a Writer, and a Reader
type DumbLite struct {
	FactomdLocation string
}

func NewDumbLite() *DumbLite {
	d := new(DumbLite)
	d.FactomdLocation = "localhost:8088"
	factom.SetFactomdServer(d.FactomdLocation)

	return d
}

//
// Writer Interface Methods
//

func (d *DumbLite) SubmitEntry(e factom.Entry, ec factom.ECAddress) (comId string, eHash string, err error) {
	comId, err = factom.CommitEntry(&e, &ec)
	if err != nil {
		return "", "", err
	}

	eHash, err = factom.RevealEntry(&e)
	return
}

func (d *DumbLite) SubmitChain(c factom.Chain, ec factom.ECAddress) (comId string, chainID string, err error) {
	comId, err = factom.CommitChain(&c, &ec)
	if err != nil {
		return "", "", err
	}

	chainID, err = factom.RevealChain(&c)

	return
}

//
// Reader Interface Methods
//

func (d *DumbLite) GetAllChainEntries(chainID primitives.Hash) ([]*factom.Entry, error) {
	return factom.GetAllChainEntries(chainID.String())
}

func (d *DumbLite) GetFirstEntry(chainID primitives.Hash) (*factom.Entry, error) {
	return factom.GetFirstEntry(chainID.String())
}

func (d *DumbLite) GetEntry(entryHash primitives.Hash) (*factom.Entry, error) {
	return factom.GetEntry(entryHash.String())
}

// GetReadyHeight indicates which height is ready to be grabbed and parsed up to.
// Calling any read function above this height will either fail, or return incomplete
// data sets.
func (d *DumbLite) GetReadyHeight() (uint32, error) {
	h, err := factom.GetHeights()
	if err != nil {
		return 0, err
	}
	return uint32(h.EntryHeight), nil
}

// GrabAllEntriesAtHeight grabs all the entries that we care about at a given height
// The EntryHolder has the entries timestamp included
func (d *DumbLite) GrabAllEntriesAtHeight(height uint32) ([]*EntryHolder, error) {
	rh, err := d.GetReadyHeight()
	if err != nil {
		return nil, err
	}
	if height > rh {
		return nil, fmt.Errorf("The height given is not ready to be grabbed and parsed. Given %d, ready up to %d", height, rh)
	}

	dblockRaw, err := factom.GetBlockByHeightRaw("d", int64(height))
	if err != nil {
		return nil, err
	}

	data, err := dblockRaw.DBlock.MarshalJSON()
	if err != nil {
		return nil, err
	}

	dblock := new(DBlock)
	err = json.Unmarshal(data, dblock)
	if err != nil {
		return nil, err
	}

	entries := make([]*EntryHolder, 0)
	// Cycle through entry blocks
	for _, eb := range dblock.Dbentries {
		if !validChain(eb.Chainid) {
			continue // Ignore this eblock, not ours
		}

		eblock, err := factom.GetEBlock(eb.Keymr)
		if err != nil {
			return nil, err
		}

		// Cycle through entries
		ents := eblock.EntryList
		for _, e := range ents {
			entry, err := factom.GetEntry(e.EntryHash)
			if err != nil {
				return nil, err
			}

			eholder := new(EntryHolder)
			eholder.Timestamp = e.Timestamp
			eholder.Entry = entry
			eholder.Height = uint32(dblock.Header.Dbheight)

			entries = append(entries, eholder)
		}
	}

	return entries, nil
}

// Check if it is a chain we care about
func validChain(chain string) bool {
	pre := hex.EncodeToString(constants.CHAIN_PREFIX)
	if chain[:constants.CHAIN_PREFIX_LENGTH_CHECK*2] == pre[:constants.CHAIN_PREFIX_LENGTH_CHECK*2] {
		return true
	}
	return false
}
