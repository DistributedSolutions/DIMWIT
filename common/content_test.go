package common_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	. "github.com/DistributedSolutions/DIMWIT/common"
	"github.com/DistributedSolutions/DIMWIT/common/primitives"
	"github.com/DistributedSolutions/DIMWIT/common/primitives/random"
)

var _ = fmt.Sprintf("")

func TestContent(t *testing.T) {
	for i := 0; i < 500; i++ {
		l := RandomNewContent()
		data, err := l.MarshalBinary()
		if err != nil {
			t.Error(err)
		}

		n := new(Content)
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

		if i > 10 {
			continue
		}

		j := new(Content)
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

	for i := 0; i < 3; i++ {
		l := RandomContentList(100)
		data, err := l.MarshalBinary()
		if err != nil {
			t.Error(err)
		}

		n := new(ContentList)
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

func TestContentIsSameAs(t *testing.T) {
	a := RandomNewContent()
	data, err := a.MarshalBinary()
	if err != nil {
		t.Fail()
	}
	b := RandomNewContent()
	err = b.UnmarshalBinary(data)
	if err != nil {
		t.Fail()
	}
	// Test IsSameAs
	for i := 0; i < 25; i++ {
		d1 := random.RandByteSliceOfSize(5000)
		d2 := random.RandByteSliceOfSize(5000)
		a.Thumbnail.SetImage(d1)
		b.Thumbnail.SetImage(d2)
		if bytes.Compare(d1, d2) != 0 && a.IsSameAs(b) {
			t.Error("Should not match")
		}
		b.Thumbnail.SetImage(d1)

		d1 = random.RandByteSliceOfSize(2)
		d2 = random.RandByteSliceOfSize(2)
		a.Type = d1[0]
		b.Type = d2[0]
		if d1[0] != d2[0] && a.IsSameAs(b) {
			t.Error("Should not match")
		}
		b.Type = d1[0]

		d1 = random.RandByteSliceOfSize(2)
		d2 = random.RandByteSliceOfSize(2)
		a.Series = d1[0]
		b.Series = d2[0]
		if d1[0] != d2[0] && a.IsSameAs(b) {
			t.Error("Should not match")
		}
		b.Series = d1[0]

		d1 = random.RandByteSliceOfSize(2)
		d2 = random.RandByteSliceOfSize(2)
		copy(a.Part[:], d1[:])
		copy(b.Part[:], d2[:])
		if bytes.Compare(d1, d2) != 0 && a.IsSameAs(b) {
			t.Error("Should not match")
		}
		copy(b.Part[:], d1[:])

		s1 := random.RandStringOfSize(a.ContentTitle.MaxLength())
		s2 := random.RandStringOfSize(b.ContentTitle.MaxLength())
		a.ContentTitle.SetString(s1)
		b.ContentTitle.SetString(s2)
		if s1 != s2 && a.IsSameAs(b) {
			t.Errorf("Should not match:%s, %s", a.ContentTitle.String(), b.ContentTitle.String())
		}
		b.ContentTitle.SetString(s1)

		s1 = random.RandStringOfSize(a.LongDescription.MaxLength())
		s2 = random.RandStringOfSize(b.LongDescription.MaxLength())
		a.LongDescription.SetString(s1)
		b.LongDescription.SetString(s2)
		if s1 != s2 && a.IsSameAs(b) {
			t.Error("Should not match")
		}
		b.LongDescription.SetString(s1)

		s1 = random.RandStringOfSize(a.ShortDescription.MaxLength())
		s2 = random.RandStringOfSize(b.ShortDescription.MaxLength())
		a.ShortDescription.SetString(s1)
		b.ShortDescription.SetString(s2)
		if s1 != s2 && a.IsSameAs(b) {
			t.Error("Should not match")
		}
		b.ShortDescription.SetString(s1)

		h1 := *primitives.RandomHash()
		h2 := *primitives.RandomHash()
		a.ContentID = h1
		b.ContentID = h2
		if !h1.IsSameAs(&h2) && a.IsSameAs(b) {
			t.Error("Should not match")
		}
		b.ContentID = h1

		h1 = *primitives.RandomHash()
		h2 = *primitives.RandomHash()
		a.RootChainID = h1
		b.RootChainID = h2
		if !h1.IsSameAs(&h2) && a.IsSameAs(b) {
			t.Error("Should not match")
		}
		b.RootChainID = h1

		i1 := *primitives.RandomInfoHash()
		i2 := *primitives.RandomInfoHash()
		a.InfoHash = i1
		b.InfoHash = i2
		if !i1.IsSameAs(&i2) && a.IsSameAs(b) {
			t.Error("Should not match")
		}
		b.InfoHash = i1

		f1 := *primitives.RandomFileList(10)
		f2 := *primitives.RandomFileList(10)
		a.ActionFiles = f1
		b.ActionFiles = f2
		if !f1.IsSameAs(&f2) && a.IsSameAs(b) {
			t.Error("Should not match")
		}
		b.ActionFiles = f1

		f1 = *primitives.RandomFileList(10)
		f2 = *primitives.RandomFileList(10)
		a.FileList = f1
		b.FileList = f2
		if !f1.IsSameAs(&f2) && a.IsSameAs(b) {
			t.Error("Should not match")
		}
		b.FileList = f1

		tr1 := *primitives.RandomTrackerList(10)
		tr2 := *primitives.RandomTrackerList(10)
		a.Trackers = tr1
		b.Trackers = tr2
		if !tr1.IsSameAs(&tr2) && a.IsSameAs(b) {
			t.Error("Should not match")
		}
		b.Trackers = tr1

		t1 := *primitives.RandomTagList(10)
		t2 := *primitives.RandomTagList(10)
		a.Tags = t1
		b.Tags = t2
		if !t1.IsSameAs(&t2) && a.IsSameAs(b) {
			t.Error("Should not match")
		}
		b.Tags = t1
	}
}

func TestBadUnmarshalContent(t *testing.T) {
	badData := []byte{}

	n := new(Content)

	_, err := n.UnmarshalBinaryData(badData)
	if err == nil {
		t.Error("Should panic or error out")
	}
}
