package objects

import (
	"fmt"

	//"github.com/DistributedSolutions/DIMWIT/constructor/entries"
	"github.com/DistributedSolutions/DIMWIT/factom-lite"
)

// ParseFactomEntry takes a factom entry, and probably parses it. Returning
// an ApplyEntry. An ApplyEntry has request functions, which the apply loop
// will honor. After giving the ApplyEntry all the info it needs, it can
// do the correct channel changes and return the new channel
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
	case "Channel Root Chain": // Root Create
		iae = NewRootChainApplyEntry()
	case "Channel Management Chain": // Manage Create
		iae = NewManageChainApplyEntry()
	case "Channel Content Chain": // Content Create
		iae = NewContentChainApplyEntry()
	case "Channel Chain": // Register Root
		iae = NewRootRegisterApplyEntry()
	case "Register Management Chain": // Register Manage
		iae = NewManageRegisterApplyEntry()
	case "Register Content Chain": // Register Content
		iae = NewContentRegisterApplyEntry()
	case "Channel Management Metadata Main":
	case "Content Signing Key":
		iae = NewContentSigningKeyApplyEntry()
	case "Content Link": // Hyperlink
		iae = NewContentLinkApplyEntry()
	case "Content Chain": // Need to stich entries in the chain.
		// We actually process Content Chains by the "Content Link", so we can
		// toss these
		iae = NewBitBucketApplyEntry()
	default:
		// Toss the entry, I have no clue what it is, do you?
		iae = NewBitBucketApplyEntry()
	}

	err = iae.ParseFactomEntry(e)
	if err != nil {
		return nil, err
	}

	return iae, nil
}
