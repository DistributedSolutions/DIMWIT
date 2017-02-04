package primitives_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"

	. "github.com/DistributedSolutions/DIMWIT/common/primitives"
	"github.com/DistributedSolutions/DIMWIT/common/primitives/random"
	//"github.com/DistributedSolutions/DIMWIT/common/constants"
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
	ta := new(Tag)
	properties := gopter.NewProperties(nil)

	properties.Property("Setting tags and testing IsSameAs", prop.ForAll(
		func(t1 string, t2 string) bool {
			if len(t1) > ta.MaxLength() {
				t1 = t1[:ta.MaxLength()]
			}
			if len(t2) > ta.MaxLength() {
				t2 = t2[:ta.MaxLength()]
			}

			a, err1 := NewTag(t1)
			b, err2 := NewTag(t2)
			same := false
			if t1 == t2 {
				same = true
			}

			return err1 == nil && err2 == nil && a.IsSameAs(b) == same
		},
		gen.AnyString(),
		gen.AnyString(),
	))

	//properties.Run(gopter.ConsoleReporter(true))
	properties.TestingRun(t)

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
			t.Fail()
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
	a := RandomTagList(random.RandomUInt32Between(20, 100))

	b := RandomTagList(random.RandomUInt32Between(20, 100))

	if a.IsSameAs(b) {
		t.Fail()
	}

	for i := 0; i < 10; i++ {
		a.AddTagByName(fmt.Sprintf("%d", i))
		b.AddTagByName(fmt.Sprintf("%d", i))
	}

	if len(a.GetTags()) < 6 {
		t.Error("Length not long enough")
		t.Fail()
	}
	a.SetTagTo(5, "one")
	if a.GetTags()[5].String() != "one" {
		t.Errorf("Should be 'one', found %s", a.GetTags()[5].String())
	}

	b.SetTagTo(5, "two")
	if b.GetTags()[5].String() != "two" {
		t.Errorf("Should be 'two', found %s", b.GetTags()[5].String())
	}

	if i, h := b.Has("two"); i != 5 || h == false {
		t.Error("Not found, but should")
	}

	c := NewTagList(20)

	err := c.AddTagByName("Another")
	if err != nil {
		t.Error(err)
	}

	if _, h := c.Has("Another"); h == false {
		t.Error("Not found, but should")
	}

	err = c.RemoveTagByName("Another")
	if err != nil {
		t.Error(err)
	}

	if _, h := c.Has("Another"); h != false {
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
