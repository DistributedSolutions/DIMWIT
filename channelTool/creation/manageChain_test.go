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

func TestManageChain(t *testing.T) {
	p := make([]primitives.PublicKey, 3)
	for i := range p {
		p[i] = *primitives.RandomPublicKey()
	}

	rc := primitives.RandomHash()

	mc := new(ManageChain)
	sec, _ := primitives.RandomPrivateKey()
	mc.CreateManagementChain(*rc, *sec)

	data, err := hex.DecodeString(mc.Create.Chain.FirstEntry.ChainID)
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	if bytes.Compare(data[:constants.CHAIN_PREFIX_LENGTH_CHECK], constants.CHAIN_PREFIX[:constants.CHAIN_PREFIX_LENGTH_CHECK]) != 0 {
		t.Error("Invalid chainID")
	}
}

func TestMetaData(t *testing.T) {
	for i := 0; i < 1000; i++ {
		m := RandomManageChainMetaData()
		data, err := m.MarshalBinary()
		if err != nil {
			t.Error(err)
		}

		x := NewManageChainMetaData()
		newData, err := x.UnmarshalBinaryData(data)
		if err != nil {
			t.Error(err)
		}

		if len(newData) != 0 {
			t.Errorf("Newdata should have 0 bytes, it has %d", len(newData))
		}

		if !m.IsSameAs(x) {
			t.Error("Should be same")
		}
	}

	m := NewManageChainMetaData()
	data, err := m.MarshalBinary()
	if err != nil {
		t.Error(err)
	}

	x := NewManageChainMetaData()
	newData, err := x.UnmarshalBinaryData(data)
	if err != nil {
		t.Error(err)
	}

	if len(newData) != 0 {
		t.Errorf("Newdata should have 0 bytes, it has %d", len(newData))
	}

	if !m.IsSameAs(x) {
		t.Error("Should be same")
	}

}

func TestManageEntries(t *testing.T) {
	for i := 0; i < 1; i++ {
		fmt.Println(".")

		rc := primitives.RandomHash()
		mc := primitives.RandomHash()

		cc := new(ManageChain)
		sec, _ := primitives.RandomPrivateKey()

		//contentType byte, contentData ContentChainContent, root primitives.Hash, contentSignKey primitives.PrivateKey
		cD := new(ManageChainMetaData)
		if i > 10 {
			cD = RandomHugeManageChainMetaData()
		} else {
			cD = RandomHugeManageChainMetaData()
		}

		conData, err := cD.MarshalBinary()
		if err != nil {
			t.Error(err)
			t.Fail()
		}

		err = cc.CreateMetadata(cD, *rc, *mc, *sec)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}

		chainID := mc.String()
		for _, e := range cc.MetaData.Entries {
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

		if len(conData) > constants.ENTRY_MAX_SIZE && len(cc.MetaData.Entries) < 1 {
			t.Errorf("Should have more supporting entries. ConData length is %d, entrycount is %d", len(conData), len(cc.MetaData.Entries))
		}

		if len(cc.MetaData.Entries) > 0 {
			if ExIDLength(cc.MetaData.MainEntry.ExtIDs)+len(cc.MetaData.MainEntry.Content) != constants.ENTRY_MAX_SIZE {
				t.Errorf("Fist entry length is not %d bytes. It has more entries, and it's length is %d", constants.ENTRY_MAX_SIZE, ExIDLength(cc.MetaData.MainEntry.ExtIDs)+len(cc.MetaData.MainEntry.Content))
			}
		}

		for _, e := range cc.MetaData.Entries {
			if e.ChainID != mc.String() {
				t.Errorf("Bad chainID in entry")
			}
		}
	}
}
