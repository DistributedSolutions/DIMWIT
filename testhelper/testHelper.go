package testhelper

import (
	"fmt"
	"time"

	"github.com/DistributedSolutions/DIMWIT/channelTool"
	"github.com/DistributedSolutions/DIMWIT/channelTool/creation"
	"github.com/DistributedSolutions/DIMWIT/common"
	"github.com/DistributedSolutions/DIMWIT/common/primitives"
	"github.com/DistributedSolutions/DIMWIT/constructor"
	"github.com/DistributedSolutions/DIMWIT/database"
	"github.com/DistributedSolutions/DIMWIT/factom-lite"
	"github.com/FactomProject/factom"
)

var _ = fmt.Sprintf("")
var _ = time.Second

func AddChannelsToClient(fake lite.FactomLite, amt int, small bool) ([]common.Channel, error) {
	ec := lite.GetECAddress()
	chanList := make([]common.Channel, 0)
	inc := new(factom.Entry)
	inc.Content = []byte("Increment")
	for i := 0; i < amt; i++ {
		fake.SubmitEntry(*inc, *ec)

		var ch *common.Channel
		if small {
			ch = common.RandomNewSmallChannel()
		} else {
			ch = common.RandomNewChannel()
		}
		auth, err := channelTool.MakeNewAuthChannel(ch, ec)
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
			fake.SubmitEntry(*inc, *ec)
			_, _, err := fake.SubmitChain(*c, *ec)
			if err != nil {
				return nil, err
			}
		}

		eHashes := make([]string, 0)

		//"Content Signing Key"
		choose := func(match string, incr bool) {
			for i, e := range entries {
				if string(e.ExtIDs[1]) == match {
					if incr {
						fake.SubmitEntry(*inc, *ec) // Increment height
					}
					_, ehash, _ := fake.SubmitEntry(*e, *ec)
					entries[i] = entries[len(entries)-1]
					entries = entries[:len(entries)-1]
					eHashes = append(eHashes, ehash)
				}
			}
		}

		choose("Content Signing Key", true)
		choose("Register Management Chain", true)
		choose("Register Content Chain", true)
		//fake.SubmitEntry(*inc, *ec)

		for _, e := range entries {
			fake.SubmitEntry(*inc, *ec)
			_, ehash, err := fake.SubmitEntry(*e, *ec)
			if err != nil {
				return nil, err
			}
			eHashes = append(eHashes, ehash)
		}

		//fake.SubmitEntry(*inc, *ec) // Increment height

		for _, h := range eHashes {
			hash, _ := primitives.HexToHash(h)
			_, err := fake.GetEntry(*hash)
			if err != nil {
				return nil, err
			}
		}
		//fake.SubmitEntry(*inc, *ec) // Increment height
		//fake.SubmitEntry(*inc, *ec) // Increment height
		fake.SubmitEntry(*inc, *ec)
	}

	fake.SubmitEntry(*inc, *ec)
	fake.SubmitEntry(*inc, *ec)
	return chanList, nil
}

func AddChannelsFromFileToClient(fake lite.FactomLite, channels *common.ChannelList, small bool) error {
	ec := lite.GetECAddress()
	chanList := make([]common.Channel, 0)
	inc := new(factom.Entry)
	inc.Content = []byte("Increment")
	for _, ch := range channels.List {
		fake.SubmitEntry(*inc, *ec)

		auth, err := channelTool.MakeNewAuthChannel(&ch, ec)
		if err != nil {
			return err
		}

		chanList = append(chanList, auth.Channel)

		chains, err := auth.ReturnFactomChains()
		if err != nil {
			return err
		}

		entries, err := auth.ReturnFactomEntries()
		if err != nil {
			return err
		}

		for _, c := range chains {
			fake.SubmitEntry(*inc, *ec)
			_, _, err := fake.SubmitChain(*c, *ec)
			if err != nil {
				return err
			}
		}

		eHashes := make([]string, 0)

		//"Content Signing Key"
		choose := func(match string, incr bool) {
			for i, e := range entries {
				if string(e.ExtIDs[1]) == match {
					if incr {
						fake.SubmitEntry(*inc, *ec) // Increment height
					}
					_, ehash, _ := fake.SubmitEntry(*e, *ec)
					entries[i] = entries[len(entries)-1]
					entries = entries[:len(entries)-1]
					eHashes = append(eHashes, ehash)
				}
			}
		}

		choose("Content Signing Key", true)
		choose("Register Management Chain", true)
		choose("Register Content Chain", true)
		//fake.SubmitEntry(*inc, *ec)

		for _, e := range entries {
			fake.SubmitEntry(*inc, *ec)
			_, ehash, err := fake.SubmitEntry(*e, *ec)
			if err != nil {
				return err
			}
			eHashes = append(eHashes, ehash)
		}

		//fake.SubmitEntry(*inc, *ec) // Increment height

		for _, h := range eHashes {
			hash, _ := primitives.HexToHash(h)
			_, err := fake.GetEntry(*hash)
			if err != nil {
				return err
			}
		}
		//fake.SubmitEntry(*inc, *ec) // Increment height
		//fake.SubmitEntry(*inc, *ec) // Increment height
		fake.SubmitEntry(*inc, *ec)
	}

	fake.SubmitEntry(*inc, *ec)
	fake.SubmitEntry(*inc, *ec)
	return nil
}

func IncrementFakeHeight(fake lite.FactomLite) (uint32, error) {
	inc := new(factom.Entry)
	inc.Content = []byte("Increment")
	ec := lite.GetECAddress()
	fake.SubmitEntry(*inc, *ec)
	return fake.GetReadyHeight()

}

func PopulateFakeClient(small bool, amt int) (lite.FactomLite, []common.Channel, error) {
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
