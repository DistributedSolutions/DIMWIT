package common_test

import (
	"fmt"
	"testing"

	. "github.com/DistributedSolutions/DIMWIT/common"
	"github.com/DistributedSolutions/DIMWIT/common/primitives/random"
)

var _ = fmt.Sprintf("")

func TestManyPlayList(t *testing.T) {
	for i := 0; i < 1000; i++ {
		max := random.RandomUInt32Between(0, 25)
		l := RandomManyPlayList(max)
		data, err := l.MarshalBinary()
		if err != nil {
			t.Error(err)
		}

		n := new(ManyPlayList)
		newData, err := n.UnmarshalBinaryData(data)
		if err != nil {
			t.Error(err)
		}
		if !n.IsSameAs(l) {
			t.Error("Should match.")
		}

		if len(newData) != 0 {
			t.Error("Failed, should have no bytes left")
		}

		if l.Empty() && len(l.GetPlaylists()) != 0 {
			t.Error("Should not be empty")
		}
	}
}

func TestPlayList(t *testing.T) {
	for i := 0; i < 1000; i++ {
		max := random.RandomUInt32Between(0, 25)
		l := RandomSinglePlayList(max)
		data, err := l.MarshalBinary()
		if err != nil {
			t.Error(err)
		}

		n := new(SinglePlayList)
		newData, err := n.UnmarshalBinaryData(data)
		if err != nil {
			t.Error(err)
		}
		if !n.IsSameAs(l) {
			t.Error("Should match.")
		}

		if len(newData) != 0 {
			t.Error("Failed, should have no bytes left")
		}
	}
}

func TestBadUnmarshalManyPlayList(t *testing.T) {
	badData := []byte{}

	n := new(ManyPlayList)

	_, err := n.UnmarshalBinaryData(badData)
	if err == nil {
		t.Error("Should panic or error out")
	}
}
