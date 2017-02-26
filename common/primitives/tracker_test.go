package primitives_test

import (
	"encoding/json"
	"fmt"
	"testing"

	. "github.com/DistributedSolutions/DIMWIT/common/primitives"
	"github.com/DistributedSolutions/DIMWIT/common/primitives/random"
)

var _ = fmt.Sprintf("")

func TestSingleTracker(t *testing.T) {
	for i := 0; i < 1000; i++ {
		l := RandomTracker()
		data, err := l.MarshalBinary()
		if err != nil {
			t.Error(err)
		}

		n := new(Tracker)
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

func TestTrackerList(t *testing.T) {
	for i := 0; i < 100; i++ {
		l := RandomTrackerList(random.RandomUInt32Between(0, 100))
		data, err := l.MarshalBinary()
		if err != nil {
			t.Error(i, err)
		}

		n := new(TrackerList)
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

		if i > 10 {
			continue
		}

		j := new(TrackerList)
		jdata, err := json.Marshal(l)
		if err != nil {
			t.Error(err)
		}

		err = json.Unmarshal(jdata, j)
		if err != nil {
			t.Error(err)
		}

		if !n.IsSameAs(j) {
			t.Error("[JsonMarshal] Should match.")
		}
	}
}

func TestDiffTrackerList(t *testing.T) {
	same := 0
	for i := 0; i < 1000; i++ {
		a := RandomTrackerList(random.RandomUInt32Between(0, 100))
		b := RandomTrackerList(random.RandomUInt32Between(0, 100))
		if a.IsSameAs(b) {
			same++
		}
	}
	if same > 15 {
		t.Error("More than 15 are the same, it is totally random, so it is likely the IsSameAs is broken.")
	}
}

func TestEmptyTracker(t *testing.T) {
	s := new(TrackerList)
	if !s.Empty() {
		t.Error("Should be empty")
	}

	tl := new(Tracker)
	if !tl.Empty() {
		t.Error("Should be empty")
	}
}

func TestBadUnmarshalTracker(t *testing.T) {
	badData := []byte{}

	n := new(TrackerList)

	_, err := n.UnmarshalBinaryData(badData)
	if err == nil {
		t.Error("Should panic or error out")
	}

	s := new(Tracker)
	_, err = s.UnmarshalBinaryData(badData)
	if err == nil {
		t.Error("Should panic or error out")
	}
}
