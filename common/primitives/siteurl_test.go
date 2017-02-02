package primitives_test

import (
	"fmt"
	"testing"

	. "github.com/DistributedSolutions/DIMWIT/common/primitives"
)

var _ = fmt.Sprintf("")

func TestSiteUrl(t *testing.T) {
	for i := 0; i < 1000; i++ {
		l := RandomSiteURL()
		data, err := l.MarshalBinary()
		if err != nil {
			t.Error(err)
		}

		n := new(SiteURL)
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

func TestPSiteURLDiff(t *testing.T) {
	same := 0
	for i := 0; i < 1000; i++ {
		a := RandomSiteURL()
		b := RandomSiteURL()
		if a.IsSameAs(b) {
			same++
		}
	}
	if same > 15 {
		t.Error("More than 15 are the same, it is totally random, so it is likely the IsSameAs is broken.")
	}
}

func TestBadUnmarshalSURL(t *testing.T) {
	badData := []byte{}

	n := new(SiteURL)

	_, err := n.UnmarshalBinaryData(badData)
	if err == nil {
		t.Error("Should panic or error out")
	}
}
