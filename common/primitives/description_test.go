package primitives_test

import (
	"fmt"
	"testing"

	. "github.com/DistributedSolutions/DIMWIT/common/primitives"
	"github.com/DistributedSolutions/DIMWIT/common/primitives/random"
)

var _ = fmt.Sprintf("")

func TestLongDesc(t *testing.T) {
	d, err := NewLongDescription("hello")
	if err != nil {
		t.Error(err)
	}

	if d.String() != "hello" {
		t.Error("String was not set")
	}

	var _ = d
}

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

func RandomLongDescription() *LongDescription {
	l, _ := NewLongDescription("")
	l.SetString(random.RandStringOfSize(l.MaxLength()))

	return l
}

func RandomShortDescription() *ShortDescription {
	l, _ := NewShortDescription("")
	l.SetString(random.RandStringOfSize(l.MaxLength()))

	return l
}
