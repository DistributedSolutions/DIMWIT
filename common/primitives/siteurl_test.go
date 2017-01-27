package primitives_test

import (
	"fmt"
	"testing"

	. "github.com/DistributedSolutions/DIMWIT/common/primitives"
	"github.com/DistributedSolutions/DIMWIT/common/primitives/random"
)

var _ = fmt.Sprintf("")

func TestSiteUrl(t *testing.T) {
	for i := 0; i < 1000; i++ {
		l := RandomSitUrl()
		data, err := l.MarshalBinary()
		if err != nil {
			t.Error(err)
		}

		n := new(SiteURL)
		err = n.UnmarshalBinary(data)
		if err != nil {
			t.Error(err)
		}
		if !n.IsSameAs(l) {
			t.Error("Should match.")
		}
	}
}

func RandomSitUrl() *SiteURL {
	l, _ := NewURL("")
	l.SetString(random.RandStringOfSize(l.MaxLength()))

	return l
}
