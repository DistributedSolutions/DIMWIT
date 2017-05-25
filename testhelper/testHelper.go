package testhelper

import (
	"fmt"
	"time"

	"github.com/DistributedSolutions/DIMWIT/common"
	//	"github.com/DistributedSolutions/DIMWIT/common/primitives"
	"github.com/DistributedSolutions/DIMWIT/constructor"
	"github.com/DistributedSolutions/DIMWIT/database"
	"github.com/DistributedSolutions/DIMWIT/factom-lite"
	"github.com/DistributedSolutions/DIMWIT/writeHelper"
	"github.com/FactomProject/factom"
)

var _ = fmt.Sprintf("")
var _ = time.Second

func AddChannelsToClient(fake lite.FactomLite, amt int, small bool) ([]common.Channel, error) {
	ec := lite.GetECAddress()
	chanList := make([]common.Channel, 0)
	inc := new(factom.Entry)
	inc.Content = []byte("Increment")

	db := database.NewMapDB()
	con, err := constructor.NewContructor(db, new(constructor.FakeSqlWriter))
	if err != nil {
		return nil, err
	}
	w, err := writeHelper.NewWriterHelper(con, fake)
	if err != nil {
		return nil, err
	}

	for i := 0; i < amt; i++ {
		fake.SubmitEntry(*inc, *ec)

		var ch *common.Channel
		if small {
			ch = common.RandomNewSmallChannel()
		} else {
			ch = common.RandomNewChannel()
		}

		err := w.MakeNewAuthChannel(ch)
		if err != nil {
			return nil, err
		}
		chanList = append(chanList, *ch)
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

	db := database.NewMapDB()
	con, err := constructor.NewContructor(db, new(constructor.FakeSqlWriter))
	if err != nil {
		return err
	}
	w, err := writeHelper.NewWriterHelper(con, fake)
	if err != nil {
		return err
	}
	for _, ch := range channels.List {
		fake.SubmitEntry(*inc, *ec)

		chanList = append(chanList, ch)
		err := w.MakeNewAuthChannel(&ch)
		if err != nil {
			return err
		}
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
	fake := lite.NewMapFakeDumbLite()
	// m := creation.NewMasterChain()
	//ec := lite.GetECAddress()
	//	fake.SubmitChain(*m.Chain, *ec)

	chanList, err := AddChannelsToClient(fake, 5, small)
	return fake, chanList, err
}

func PopulateLevel2Cache(fake lite.FactomLite) (*constructor.Constructor, database.IDatabase, error) {
	// Starts SQL InterfaceDB
	sql, err := constructor.NewSqlWriter()
	if err != nil {
		return nil, nil, err
	}

	db := database.NewMapDB()
	con, err := constructor.NewContructor(db, sql)
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
