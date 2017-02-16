package objects

import (
	"fmt"

	//"github.com/DistributedSolutions/DIMWIT/constructor/entries"
	"github.com/DistributedSolutions/DIMWIT/factom-lite"
)

// ParseFactomEntry takes a factom entry, and probably parses it. Returning
// an ApplyEntry and the channel it needs to interact with
func ParseFactomEntry(e *lite.EntryHolder) (iae IApplyEntry, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("[ParseFactomEntry] A panic has occurred while parsing a factom entry: %s", r)
			return
		}
	}()

	if len(e.Entry.ExtIDs) < 2 {
		return nil, fmt.Errorf("ExternalID length is less than 2. This is too short")
	}

	// This is how we designate the entry type
	// and parse appropriately
	switch string(e.Entry.ExtIDs[2]) {
	case "Master Chain":
		iae = NewMasterChainApplyEntry()
	case "Channel Chain": // Ents in MasterChain
	case "Channel Root Chain":
		iae = NewRootChainApplyEntry()
	case "Content Signing Key":
	case "Register Management Chain": // Need to stich entries in the chain.
	case "Register Content Chain":
	case "Channel Management Chain":
	case "Channel Content Chain":
	case "Content Link":
	case "Content Chain": // Need to stich entries in the chain.
	}

	err = iae.ParseFactomEntry(e)
	if err != nil {
		return nil, err
	}

	return MasterChainApplyEntry, nil
}
