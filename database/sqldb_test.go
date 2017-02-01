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
	if constants.TRAVIS_RUN {
		return // Until this works for travis
	}
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
	if constants.TRAVIS_RUN {
		return // Until this works for travis
	}
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
	if constants.TRAVIS_RUN {
		return // Until this works for travis
	}
	c := common.RandomNewChannel()
	err := AddChannel(c)
	if err != nil {
		t.Error(err)
	}
	CloseDB()
}
