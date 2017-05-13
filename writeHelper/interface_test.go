package writeHelper_test

import (
	"testing"

	. "github.com/DistributedSolutions/DIMWIT/writeHelper"
)

func TestFakeWriterHelper(t *testing.T) {
	f := new(FakeWriterHelper)

	f.VerifyChannel(nil)
	f.InitiateChannel(nil)
	f.UpdateChannel(nil)
	f.DeleteChannel(nil)
	f.VerifyContent(nil)
	f.AddContent(nil)
	f.DeleteContent(nil)
}
