package database_test

import (
	"fmt"
	"github.com/DistributedSolutions/DIMWIT/common"
	"github.com/DistributedSolutions/DIMWIT/common/constants"
	"testing"

	. "github.com/DistributedSolutions/DIMWIT/database"
)

var _ = fmt.Sprintf("")

func TestCreateDB(t *testing.T) {
	err := DeleteDB(constants.SQL_DB)
	if err != nil {
		t.Error(err)
	}
	err = CreateDB(constants.SQL_DB, CREATE_TABLE)
	if err != nil {
		t.Error(err)
	}
}

func TestAddTags(t *testing.T) {
	err := DeleteTags()
	if err != nil {
		t.Error(err)
	}

	err = AddTags()
	if err != nil {
		t.Error(err)
	}
}

func TestAddChannel(t *testing.T) {
	c := common.RandomNewChannel()
	if len(c.Tags.GetTags()) < 1 {
		t.Error(fmt.Printf("Error tags count for random is bad :( [%d]\n", len(c.Tags.GetTags())))
	}
	err := AddChannel(c)
	if err != nil {
		t.Error(err)
	}
	CloseDB()
}
