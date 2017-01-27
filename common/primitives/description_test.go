package primitives_test

import (
	"fmt"
	"testing"

	. "github.com/DistributedSolutions/DIMWIT/common/primitives"
)

var _ = fmt.Sprintf("")

func TestDescriptions(t *testing.T) {
	for i := 0; i < 1000; i++ {
		l := RandomLongDescription()
		data, err := l.MarshalBinary()
		if err != nil {
			t.Error(err)
		}

		n := new(LongDescription)
		err = n.UnmarshalBinary(data)
		if err != nil {
			t.Error(err)
		}
		if !n.IsSameAs(l) {
			t.Error("Should match.")
		}
	}

	for i := 0; i < 1000; i++ {
		s := RandomShortDescription()
		data, err := s.MarshalBinary()
		if err != nil {
			t.Error(err)
		}

		n := new(ShortDescription)
		err = n.UnmarshalBinary(data)
		if err != nil {
			t.Error(err)
		}
		if !n.IsSameAs(s) {
			t.Error("Should match.")
		}
	}
}
