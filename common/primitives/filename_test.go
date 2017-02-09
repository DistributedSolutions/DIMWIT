package primitives_test

import (
	"fmt"
	"testing"

	. "github.com/DistributedSolutions/DIMWIT/common/primitives"
	"github.com/DistributedSolutions/DIMWIT/common/primitives/random"
)

var _ = fmt.Sprintf("")

func TestSingleFile(t *testing.T) {
	for i := 0; i < 1000; i++ {
		l := RandomFile()
		data, err := l.MarshalBinary()
		if err != nil {
			t.Error(err)
		}

		n := new(File)
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

func TestFileList(t *testing.T) {
	for i := 0; i < 100; i++ {
		l := RandomFileList(random.RandomUInt32Between(0, 100))
		data, err := l.MarshalBinary()
		if err != nil {
			t.Error(i, err)
		}

		n := new(FileList)
		newData, err := n.UnmarshalBinaryData(data)
		if err != nil {
			t.Error(i, err)
		}

		if !n.IsSameAs(l) {
			t.Error("Should match.")
		}
		if len(newData) != 0 {
			t.Error("Failed, should have no bytes left")
		}

		if l.Empty() && len(l.GetFiles()) != 0 {
			t.Error("Should not be empty")
		}
	}
}

func TestDiffFileList(t *testing.T) {
	same := 0
	for i := 0; i < 1000; i++ {
		a := RandomFileList(random.RandomUInt32Between(0, 1000))
		b := RandomFileList(random.RandomUInt32Between(0, 1000))
		if a.IsSameAs(b) {
			same++
		}
	}
	if same > 15 {
		t.Error("More than 15 are the same, it is totally random, so it is likely the IsSameAs is broken.")
	}

}

func TestEmptyFiles(t *testing.T) {
	af := new(FileList)
	if !af.Empty() {
		t.Error("Should be empty")
	}

	f := new(File)
	if !f.Empty() {
		t.Error("Should be empty")
	}

}

func TestBadUnmarshalFile(t *testing.T) {
	badData := []byte{}

	n := new(File)

	_, err := n.UnmarshalBinaryData(badData)
	if err == nil {
		t.Error("Should panic or error out")
	}

	s := new(FileList)
	_, err = s.UnmarshalBinaryData(badData)
	if err == nil {
		t.Error("Should panic or error out")
	}
}
