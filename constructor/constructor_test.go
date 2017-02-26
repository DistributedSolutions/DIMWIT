package constructor_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/DistributedSolutions/DIMWIT/channelTool"
	"github.com/DistributedSolutions/DIMWIT/channelTool/creation"
	"github.com/DistributedSolutions/DIMWIT/common"
	"github.com/DistributedSolutions/DIMWIT/common/constants"
	"github.com/DistributedSolutions/DIMWIT/common/primitives"
	. "github.com/DistributedSolutions/DIMWIT/constructor"
	"github.com/DistributedSolutions/DIMWIT/database"
	"github.com/DistributedSolutions/DIMWIT/factom-lite"
)

var _ = fmt.Sprintf("")
var _ = time.Second

func TestConstructor(t *testing.T) {
	constants.CHECK_FACTOM_FOR_UPDATES = time.Millisecond * 100
	constants.CHAIN_PREFIX_LENGTH_CHECK = 1
	fake := lite.NewFakeDumbLite()
	m := creation.NewMasterChain()
	ec := lite.GetECAddress()
	fake.SubmitChain(*m.Chain, *ec)
	//fake := lite.NewDumbLite()

	db := database.NewMapDB()
	con, err := NewContructor(db)
	if err != nil {
		t.Error(err)
	}
	con.SetReader(fake)
	authList := make([]channelTool.AuthChannel, 0)
	totalStuff := 0

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
			totalStuff++
			_, _, err := fake.SubmitChain(*c, *ec)
			if err != nil {
				t.Error(err)
			}
		}

		eHashes := make([]string, 0)
		for _, e := range entries {
			totalStuff++
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
		authList = append(authList, *auth)

		// Jesse Test Assertions

		// END JESSE ASSERTIONS
	}
	//fmt.Println(con.CompletedHeight)
	go con.StartConstructor()
	time.Sleep(1 * time.Second)
	max, _ := con.Reader.GetReadyHeight()
	for con.CompletedHeight < max-1 {
		//fmt.Println(con.CompletedHeight, max-1, totalStuff-1)
		time.Sleep(200 * time.Millisecond)
		// fmt.Println(con.CompletedHeight, max)
	}
	//fmt.Println(con.CompletedHeight, max-1, totalStuff-1)

	for _, a := range authList {
		// Constructor finished!
		cw, err := con.RetrieveChannel(a.Channel.RootChainID)
		if err != nil {
			t.Error(err)
		}

		if !cw.Channel.IsSameAs(&a.Channel) {
			t.Error("Channels not the same", cw.Channel.RootChainID.String(), a.Channel.RootChainID.String())
		}
	}

	// Close constructor
	con.Close()

}
