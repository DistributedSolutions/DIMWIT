package creation_test

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"testing"

	. "github.com/DistributedSolutions/DIMWIT/channelTool/creation"
	"github.com/DistributedSolutions/DIMWIT/common/constants"
	"github.com/DistributedSolutions/DIMWIT/common/primitives"
)

var _ = fmt.Sprintf("")

func init() {
	// So tests are quick, the normal prefix check length takes too long
	constants.CHAIN_PREFIX_LENGTH_CHECK = 1
}

func TestContentChain(t *testing.T) {
	for i := 0; i < 1000; i++ {
		rc := primitives.RandomHash()

		cc := new(ContentChain)
		sec, _ := primitives.RandomPrivateKey()

		//contentType byte, contentData ContentChainContent, root primitives.Hash, contentSignKey primitives.PrivateKey
		cD := new(ContentChainContent)
		if i > 10 {
			cD = RandomHugeContentChainContent()
		} else {
			cD = RandomContentChainContent()
		}

		conData, err := cD.MarshalBinary()
		if err != nil {
			t.Error(err)
			t.Fail()
		}

		err = cc.CreateContentChain(0x00, *cD, *rc, *sec)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}

		chainID := cc.FirstEntry.ChainID
		for _, e := range cc.Entries {
			if len(e.Content) > constants.ENTRY_MAX_SIZE || len(e.Content) == 0 {
				t.Errorf("Entry content length is bad. It is %d", len(e.Content))
			}

			if e.ChainID != chainID {
				t.Errorf("ChainID of entry is bad. Found %s, should be %s", e.ChainID, chainID)
			}

			if ExIDLength(e.ExtIDs)+len(e.Content) > constants.ENTRY_MAX_SIZE {
				t.Errorf("Entry length is too large. Entry length is %d", ExIDLength(e.ExtIDs)+len(e.Content))
			}
		}

		if len(conData) > constants.ENTRY_MAX_SIZE && len(cc.Entries) < 1 {
			t.Errorf("Should have more supporting entries. ConData length is %d, entrycount is %d", len(conData), len(cc.Entries))
		}

		if len(cc.Entries) > 0 {
			if ExIDLength(cc.FirstEntry.FirstEntry.ExtIDs)+len(cc.FirstEntry.FirstEntry.Content) != constants.ENTRY_MAX_SIZE {
				t.Errorf("%d - Fist entry length is not %d bytes. It has more entries, and it's length is %d. Header: %d, Content: %d",
					i, constants.ENTRY_MAX_SIZE, ExIDLength(cc.FirstEntry.FirstEntry.ExtIDs)+len(cc.FirstEntry.FirstEntry.Content),
					ExIDLength(cc.FirstEntry.FirstEntry.ExtIDs), len(cc.FirstEntry.FirstEntry.Content))
			}
		}

		data, err := hex.DecodeString(cc.FirstEntry.ChainID)
		if err != nil {
			t.Error(err)
			t.Fail()
		}

		if bytes.Compare(data[:constants.CHAIN_PREFIX_LENGTH_CHECK], constants.CHAIN_PREFIX[:constants.CHAIN_PREFIX_LENGTH_CHECK]) != 0 {
			t.Errorf("Invalid chainID, found %x, expected %x", data[:constants.CHAIN_PREFIX_LENGTH_CHECK], constants.CHAIN_PREFIX[:constants.CHAIN_PREFIX_LENGTH_CHECK])
		}
	}
}

func TestContentChainContent(t *testing.T) {
	for i := 0; i < 1000; i++ {
		l := RandomContentChainContent()
		data, err := l.MarshalBinary()
		if err != nil {
			t.Error(err)
		}

		n := new(ContentChainContent)
		newData, err := n.UnmarshalBinaryData(data)
		if err != nil {
			t.Error(err)
		}
		if !n.IsSameAs(l) {
			t.Error("Should match.")
		}

		if len(newData) != 0 {
			t.Error("Failed, should have no bytes left")
		}
	}
}

func TestBadUnmarshalContentChainContent(t *testing.T) {
	badData := []byte{}

	n := new(ContentChainContent)

	_, err := n.UnmarshalBinaryData(badData)
	if err == nil {
		t.Error("Should panic or error out")
	}
}
