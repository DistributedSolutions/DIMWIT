package primitives_test

import (
	"fmt"
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"

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

	s := new(ShortDescription)
	err := s.UnmarshalBinary(nil)
	if err == nil {
		t.Error("Should error")
	}

	l := new(LongDescription)
	err = l.UnmarshalBinary(nil)
	if err == nil {
		t.Error("Should error")
	}
}

func TestDiffDescription(t *testing.T) {
	a := RandomLongDescription()
	data, err := a.MarshalBinary()
	if err != nil {
		t.Error(err)
	}

	b := new(LongDescription)
	_, err = b.UnmarshalBinaryData(data)
	if err != nil {
		t.Error(err)
	}

	a.SetString("One")
	b.SetString("Two")

	if a.IsSameAs(b) {
		t.Error("Should be different")
	}
}

func TestProp(t *testing.T) {
	properties := gopter.NewProperties(nil)
	properties.Property("Seting Short Descs", prop.ForAll(
		func(t1 string) bool {
			s, err := NewShortDescription(t1)
			if len(t1) > s.MaxLength() {
				return err != nil
			} else {
				return err == nil
			}
		},
		gen.AnyString(),
	))

	properties.Property("Seting Long Descs", prop.ForAll(
		func(t1 string) bool {
			s, err := NewLongDescription(t1)
			if len(t1) > s.MaxLength() {
				return err != nil
			} else {
				return err == nil
			}
		},
		gen.AnyString(),
	))

	properties.TestingRun(t)

}

func TestEmptyDescs(t *testing.T) {
	d := new(LongDescription)
	if !d.Empty() {
		t.Error("Should be empty")
	}

	s := new(ShortDescription)
	if !s.Empty() {
		t.Error("Should be empty")
	}
}

func TestBadUnmarshalDesc(t *testing.T) {
	badData := []byte{}

	n := new(LongDescription)

	_, err := n.UnmarshalBinaryData(badData)
	if err == nil {
		t.Error("Should panic or error out")
	}

	s := new(ShortDescription)
	_, err = s.UnmarshalBinaryData(badData)
	if err == nil {
		t.Error("Should panic or error out")
	}
}
