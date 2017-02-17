package objects_test

import (
	"testing"

	"github.com/DistributedSolutions/DIMWIT/channelTool"
	"github.com/DistributedSolutions/DIMWIT/common"
	"github.com/DistributedSolutions/DIMWIT/factom-lite"
)

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

}
