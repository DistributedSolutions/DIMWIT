package database_test

import (
	"fmt"
	"github.com/DistributedSolutions/DIMWIT/common"
	"testing"

	. "github.com/DistributedSolutions/DIMWIT/database"
)

var _ = fmt.Sprintf("")

func TestCreateDB(t *testing.T) {
	err := DeleteDB()
	if err != nil {
		t.Error(err)
	}
	err = CreateDB()
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
	err := AddChannel(common.RandomNewChannel())
	if err != nil {
		t.Error(err)
	}
}
