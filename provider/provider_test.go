package provider_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/DistributedSolutions/DIMWIT/channelTool"
	"github.com/DistributedSolutions/DIMWIT/channelTool/creation"
	"github.com/DistributedSolutions/DIMWIT/common"
	"github.com/DistributedSolutions/DIMWIT/common/constants"
	"github.com/DistributedSolutions/DIMWIT/common/primitives"
	"github.com/DistributedSolutions/DIMWIT/constructor"
	"github.com/DistributedSolutions/DIMWIT/database"
	"github.com/DistributedSolutions/DIMWIT/factom-lite"
	. "github.com/DistributedSolutions/DIMWIT/provider"
)

var _ = fmt.Sprintf("")
var _ = time.Second

func TestProvider(t *testing.T) {
	fake := lite.NewFakeDumbLite()
	m := creation.NewMasterChain()
	ec := lite.GetECAddress()
	fake.SubmitChain(*m.Chain, *ec)
	//fake := lite.NewDumbLite()
	chanList := make([]common.Channel, 0)

	for i := 0; i < 5; i++ {
		ch := common.RandomNewChannel()
		auth, err := channelTool.NewAuthChannel(ch, ec)
		if err != nil {
			t.Error(err)
		}

		chanList = append(chanList, auth.Channel)

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

		for _, h := range eHashes {
			hash, _ := primitives.HexToHash(h)
			_, err := fake.GetEntry(*hash)
			if err != nil {
				t.Error(err)
			}
		}
	}

	db := database.NewMapDB()
	con, err := constructor.NewContructor(db)
	if err != nil {
		t.Error(err)
	}
	defer con.Close()

	con.SetReader(fake)
	go con.StartConstructor()

	max, _ := con.Reader.GetReadyHeight()
	for con.CompletedHeight < max-1 {
		time.Sleep(200 * time.Millisecond)
		// fmt.Println(con.CompletedHeight, max)
	}

	prov, err := NewProvider(db)
	if err != nil {
		t.Error(err)
	}

	for _, c := range chanList {
		nc, err := prov.GetChannel(c.RootChainID.String())
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		if !nc.IsSameAs(&c) {
			t.Error("Channel should be equal")
		}

		for _, content := range c.Content.GetContents() {
			newContent, err := prov.GetContent(content.ContentID.String())
			if err != nil {
				t.Error(err)
				t.FailNow()
			}
			if !newContent.IsSameAs(&content) {
				t.Error("Content should be equal")
			}
		}
	}

	_, err = prov.GetCompleteHeight()
	if err != nil {
		t.Error(err)
	}
}
