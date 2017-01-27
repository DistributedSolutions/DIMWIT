package primitives_test

import (
	"fmt"
	"testing"

	. "github.com/DistributedSolutions/DIMWIT/common/primitives"
	"github.com/DistributedSolutions/DIMWIT/common/primitives/random"
)

var _ = fmt.Sprintf("")

func TestTitle(t *testing.T) {
	for i := 0; i < 1000; i++ {
		l := RandomTitle()
		data, err := l.MarshalBinary()
		if err != nil {
			t.Error(err)
		}

		n := new(Title)
		err = n.UnmarshalBinary(data)
		if err != nil {
			t.Error(err)
		}
		if !n.IsSameAs(l) {
			t.Error("Should match.")
		}
	}
}

func RandomTitle() *Title {
	l, _ := NewTitle("")
	l.SetString(random.RandStringOfSize(l.MaxLength()))

	return l
}
