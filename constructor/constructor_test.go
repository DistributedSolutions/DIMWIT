package constructor_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/DistributedSolutions/DIMWIT/channelTool"
	"github.com/DistributedSolutions/DIMWIT/channelTool/creation"
	"github.com/DistributedSolutions/DIMWIT/common"
	. "github.com/DistributedSolutions/DIMWIT/constructor"
	"github.com/DistributedSolutions/DIMWIT/database"
	"github.com/DistributedSolutions/DIMWIT/factom-lite"
)

var _ = fmt.Sprintf("")
var _ = time.Second

func TestBitbucket(t *testing.T) {
	fake := lite.NewFakeDumbLite()
	m := creation.NewMasterChain()
	ec := lite.GetECAddress()
	fake.SubmitChain(*m.Chain, *ec)
	//fake := lite.NewDumbLite()

	for i := 0; i < 5; i++ {
		ch := common.RandomNewChannel()
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

		db := database.NewMapDB()
		con, err := NewContructor(db)
		if err != nil {
			t.Error(err)
		}

		con.SetReader(fake)
		go con.StartConstructor()

		/*for i := 0; ; i++ {
			err := con.ApplyHeight(uint32(i))
			if err != nil {
				break
			}
			fmt.Println(i)
		}*/
		//time.Sleep(10 * time.Second)

		max, _ := con.Reader.GetReadyHeight()
		for con.CompletedHeight < max-1 {
			time.Sleep(200 * time.Millisecond)
			// fmt.Println(con.CompletedHeight, max)
		}

		cw, err := con.RetrieveChannel(auth.Channel.RootChainID)
		if err != nil {
			t.Error(err)
		}

		if !cw.Channel.IsSameAs(&auth.Channel) {
			t.Error("Channels not the same", cw.Channel.RootChainID.String(), auth.Channel.RootChainID.String())
		}
		con.Close()
	}

}
