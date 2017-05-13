package constructor_test

import (
	"fmt"
	"testing"
	"time"

	//"github.com/DistributedSolutions/DIMWIT/common"
	"github.com/DistributedSolutions/DIMWIT/common/constants"
	//"github.com/DistributedSolutions/DIMWIT/common/primitives"
	. "github.com/DistributedSolutions/DIMWIT/constructor"
	"github.com/DistributedSolutions/DIMWIT/database"
	//"github.com/DistributedSolutions/DIMWIT/factom-lite"
	"github.com/DistributedSolutions/DIMWIT/testHelper"
)

var _ = fmt.Sprintf("")
var _ = time.Second

func TestConstructor(t *testing.T) {
	constants.CHECK_FACTOM_FOR_UPDATES = time.Millisecond * 100
	fake, channels, err := testhelper.PopulateFakeClient(true, 5)
	if err != nil {
		t.Error(err)
	}

	sqlW, err := NewSqlWriter()
	if err != nil {
		t.Error(err)
	}

	db := database.NewMapDB()
	con, err := NewContructor(db, sqlW)
	if err != nil {
		t.Error(err)
	}
	con.SetReader(fake)

	//fmt.Println(con.CompletedHeight)
	testhelper.IncrementFakeHeight(fake)
	testhelper.IncrementFakeHeight(fake)
	testhelper.IncrementFakeHeight(fake)
	go con.StartConstructor()
	time.Sleep(1 * time.Second)
	max, _ := con.Reader.GetReadyHeight()
	for con.CompletedHeight < max-1 {
		//fmt.Println(con.CompletedHeight, max-1, totalStuff-1)
		time.Sleep(200 * time.Millisecond)
		// fmt.Println(con.CompletedHeight, max)
	}
	//fmt.Println(con.CompletedHeight, max-1, totalStuff-1)

	for _, a := range channels {
		// Constructor finished!
		cw, err := con.RetrieveChannel(a.RootChainID)
		if err != nil {
			t.Error(err)
			continue
		}
		if cw == nil {
			t.Error("Channel not found")
			continue
		}

		if !cw.Channel.IsSameAs(&a) {
			t.Error("Channels not the same", cw.Channel.RootChainID.String(), a.RootChainID.String())
		}
	}

	// Close constructor
	con.Close()

}
