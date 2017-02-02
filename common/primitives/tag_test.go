package primitives_test

import (
	"bytes"
	"fmt"
	"testing"

	. "github.com/DistributedSolutions/DIMWIT/common/primitives"
	"github.com/DistributedSolutions/DIMWIT/common/primitives/random"
)

var _ = fmt.Sprintf("")

func TestSingleTags(t *testing.T) {
	for i := 0; i < 1000; i++ {
		l := RandomTag()

		data, err := l.MarshalBinary()
		if err != nil {
			t.Error(err)
		}

		n := new(Tag)
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

	var err error
	for i := 0; i < 1000; i++ {
		buf := new(bytes.Buffer)

		ll := make([]Tag, 0)
		ln := make([]Tag, 5)
		for c := 0; c < 5; c++ {
			l := RandomTag()
			data, err := l.MarshalBinary()
			if err != nil {
				t.Error(err)
			}
			buf.Write(data)

			ll = append(ll, *l)
		}

		newData := buf.Next(buf.Len())

		for i := range ll {
			newData, err = ln[i].UnmarshalBinaryData(newData)
			if err != nil {
				t.Error(err)
			}

			if !ln[i].IsSameAs(&ll[i]) {
				t.Error("not same")
			}
		}
	}
}

func TestTagDiff(t *testing.T) {
	a := RandomTag()
	a.SetString("one")
	b := RandomTag()
	b.SetString("two")
	if a.IsSameAs(b) {
		t.Fail()
	}

}

func TestTagList(t *testing.T) {
	for i := 0; i < 10000; i++ {
		l := RandomTagList(random.RandomUInt32Between(0, 100))
		data, err := l.MarshalBinary()
		if err != nil {
			t.Error(i, err)
		}

		n := new(TagList)
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
	}
}

func TestTagListDiff(t *testing.T) {
	a := RandomTagList(random.RandomUInt32Between(10, 100))
	a.SetTagTo(5, "one")
	b := RandomTagList(random.RandomUInt32Between(10, 100))
	b.SetTagTo(5, "two")
	if a.IsSameAs(b) {
		t.Fail()
	}

	if a.GetTags()[5] != "one" {
		t.Error("Should be 'one'")
	}

	if b.GetTags()[5] != "two" {
		t.Error("Should be 'one'")
	}

	if i, h := b.Has("two"); i != 5 || h == false {
		t.Error("Not found, but should")
	}

	err := b.AddTagByName("Another")
	if err != nil {
		t.Error(err)
	}

	if _, h := b.Has("Another"); h == false {
		t.Error("Not found, but should")
	}

	err = b.RemoveTagByName("Another")
	if err != nil {
		t.Error(err)
	}

	if _, h := b.Has("Another"); h != false {
		t.Error("Found, but shouldn't")
	}

}

func TestBadUnmarshalTag(t *testing.T) {
	badData := []byte{}

	n := new(Tag)

	_, err := n.UnmarshalBinaryData(badData)
	if err == nil {
		t.Error("Should panic or error out")
	}

	s := new(TagList)
	_, err = s.UnmarshalBinaryData(badData)
	if err == nil {
		t.Error("Should panic or error out")
	}
}
