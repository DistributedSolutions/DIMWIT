package constructor_test

import (
	"fmt"
	"testing"

	"github.com/DistributedSolutions/DIMWIT/channelTool"
	"github.com/DistributedSolutions/DIMWIT/common"
	. "github.com/DistributedSolutions/DIMWIT/constructor"
	"github.com/DistributedSolutions/DIMWIT/factom-lite"
)

var _ = fmt.Sprintf("")

func TestBitbucket(t *testing.T) {
	fake := lite.NewFakeDumbLite()
	ch := common.RandomNewChannel()
	ec := lite.GetECAddress()
	auth, err := channelTool.NewAuthChannel(ch, ec)
	if err != nil {
		t.Error(err)
	}

	chains, err := auth.ReturnFactomChains()
	if err != nil {
		t.Error(err)
	}

	entries, err := auth.ReturnFactomEntries()
	if err != nil {
		t.Error(err)
	}

	for _, c := range chains {
		_, _, err := fake.SubmitChain(*c, *ec)
		if err != nil {
			t.Error(err)
		}
	}

	eHashes := make([]string, 0)
	for _, e := range entries {
		_, ehash, err := fake.SubmitEntry(*e, *ec)
		if err != nil {
			t.Error(err)
		}
		eHashes = append(eHashes, ehash)
	}

	con, err := NewContructor("Map")
	if err != nil {
		t.Error(err)
	}

	con.SetReader(fake)
	//con.

	for i := 0; ; i++ {
		err := con.ApplyHeight(uint32(i))
		if err != nil {
			break
		}
	}

	cw, err := con.RetrieveChannel(auth.Channel.RootChainID)
	if err != nil {
		t.Error(err)
	}

	if !cw.Channel.IsSameAs(&auth.Channel) {
		t.Error("Channels not the same")
	}

}
