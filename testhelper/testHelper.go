package testhelper

import (
	"fmt"
	"time"

	"github.com/DistributedSolutions/DIMWIT/channelTool"
	"github.com/DistributedSolutions/DIMWIT/channelTool/creation"
	"github.com/DistributedSolutions/DIMWIT/common"
	"github.com/DistributedSolutions/DIMWIT/common/constants"
	"github.com/DistributedSolutions/DIMWIT/common/primitives"
	"github.com/DistributedSolutions/DIMWIT/constructor"
	"github.com/DistributedSolutions/DIMWIT/database"
	"github.com/DistributedSolutions/DIMWIT/factom-lite"
)

var _ = fmt.Sprintf("")
var _ = time.Second

func AddChannelsToClient(fake lite.FactomLite, amt int, small bool) ([]common.Channel, error) {
	ec := lite.GetECAddress()
	chanList := make([]common.Channel, 0)
	for i := 0; i < amt; i++ {
		var ch *common.Channel
		if small {
			ch = common.RandomNewSmallChannel()
		} else {
			ch = common.RandomNewChannel()
		}
		auth, err := channelTool.NewAuthChannel(ch, ec)
		if err != nil {
			return nil, err
		}

		chanList = append(chanList, auth.Channel)

		chains, err := auth.ReturnFactomChains()
		if err != nil {
			return nil, err
		}

		entries, err := auth.ReturnFactomEntries()
		if err != nil {
			return nil, err
		}

		for _, c := range chains {
			_, _, err := fake.SubmitChain(*c, *ec)
			if err != nil {
				return nil, err
			}
		}

		eHashes := make([]string, 0)
		for _, e := range entries {
			_, ehash, err := fake.SubmitEntry(*e, *ec)
			if err != nil {
				return nil, err
			}
			eHashes = append(eHashes, ehash)
		}

		for _, h := range eHashes {
			hash, _ := primitives.HexToHash(h)
			_, err := fake.GetEntry(*hash)
			if err != nil {
				return nil, err
			}
		}
	}
	return chanList, nil
}

func PopulateFakeClient(small bool, amt int) (lite.FactomLite, []common.Channel, error) {
	constants.CHAIN_PREFIX_LENGTH_CHECK = 1
	fake := lite.NewFakeDumbLite()
	m := creation.NewMasterChain()
	ec := lite.GetECAddress()
	fake.SubmitChain(*m.Chain, *ec)

	chanList, err := AddChannelsToClient(fake, 5, small)
	return fake, chanList, err
}

func PopulateLevel2Cache(fake lite.FactomLite) (*constructor.Constructor, database.IDatabase, error) {
	db := database.NewMapDB()
	con, err := constructor.NewContructor(db)
	if err != nil {
		return nil, nil, err
	}

	con.SetReader(fake)
	go con.StartConstructor()

	// Run through blockchain
	max, _ := con.Reader.GetReadyHeight()
	for con.CompletedHeight < max-1 {
		time.Sleep(200 * time.Millisecond)
		// fmt.Println(con.CompletedHeight, max)
	}

	return con, db, nil
}
