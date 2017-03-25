package provider_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/DistributedSolutions/DIMWIT/constructor"
	"github.com/DistributedSolutions/DIMWIT/database"
	. "github.com/DistributedSolutions/DIMWIT/provider"
	"github.com/DistributedSolutions/DIMWIT/testhelper"
)

var _ = fmt.Sprintf("")
var _ = time.Second

func TestProvider(t *testing.T) {
	fake, chanList, err := testhelper.PopulateFakeClient(true, 5)
	if err != nil {
		t.Error(err)
		t.FailNow()
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
	if max == 0 {
		max = 1
	}
	for con.CompletedHeight < max-1 {
		time.Sleep(200 * time.Millisecond)
		// fmt.Println(con.CompletedHeight, max)
	}

	prov, err := NewProvider(db, fake)
	if err != nil {
		t.Error(err)
	}
	// defer prov.Close()

	for _, c := range chanList {
		nc, err := prov.GetChannel(c.RootChainID.String())
		if err != nil {
			t.Error("Not found,", err.Error())
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
