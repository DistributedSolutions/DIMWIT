package constructor_test

import (
	"testing"

	"github.com/DistributedSolutions/DIMWIT/common"
	. "github.com/DistributedSolutions/DIMWIT/constructor"
)

func TestFakeSQLWriter(t *testing.T) {
	s := new(FakeSqlWriter)
	s.AddChannelArr([]common.Channel{}, 0)
	s.FlushTempPlaylists(0)
	s.DeleteDBChannels()
	s.DeleteDB()
	s.Close()
}
