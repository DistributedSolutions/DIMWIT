package common_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"testing"

	. "github.com/DistributedSolutions/DIMWIT/common"
	"github.com/DistributedSolutions/DIMWIT/common/primitives"
	"github.com/DistributedSolutions/DIMWIT/common/primitives/random"
)

var _ = fmt.Sprintf("")
var _ = log.Prefix()
var _ = ioutil.Discard

func TestChannel(t *testing.T) {
	for i := 0; i < 250; i++ {
		l := RandomNewChannel()
		data, err := l.MarshalBinary()
		if err != nil {
			t.Error(err)
		}

		n := NewChannel()
		newData, err := n.UnmarshalBinaryData(data)
		if err != nil {
			t.Error(err)
		}
		if !n.IsSameAs(l) {
			t.Error("[BytesMarshal] Should match.")
		}

		if len(newData) != 0 {
			t.Error("Failed, should have no bytes left")
		}

		if i > 10 {
			continue
		}
		j := new(Channel)
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

	a := RandomNewChannel()
	data, err := a.MarshalBinary()
	if err != nil {
		t.Error(err)
	}

	_, _ = a.CustomMarshalJSON()

	b := NewChannel()
	_, err = b.UnmarshalBinaryData(data)
	if err != nil {
		t.Error(err)
	}

	if !a.IsSameAs(b) {
		t.Error("Should match.")
	}

	// Test IsSameAs
	for i := 0; i < 100; i++ {
		s1 := random.RandStringOfSize(a.Website.MaxLength())
		s2 := random.RandStringOfSize(b.Website.MaxLength())
		a.Website.SetString(s1)
		b.Website.SetString(s2)
		if s1 != s2 && a.IsSameAs(b) {
			t.Error("Should not match")
		}
		b.Website.SetString(s1)

		d1 := random.RandByteSliceOfSize(5000)
		d2 := random.RandByteSliceOfSize(5000)
		a.Banner.SetImage(d1)
		b.Banner.SetImage(d2)
		if bytes.Compare(d1, d2) != 0 && a.IsSameAs(b) {
			t.Error("Should not match")
		}
		b.Banner.SetImage(d1)

		d1 = random.RandByteSliceOfSize(5000)
		d2 = random.RandByteSliceOfSize(5000)
		a.Thumbnail.SetImage(d1)
		b.Thumbnail.SetImage(d2)
		if bytes.Compare(d1, d2) != 0 && a.IsSameAs(b) {
			t.Error("Should not match")
		}
		b.Thumbnail.SetImage(d1)

		s1 = random.RandStringOfSize(a.ChannelTitle.MaxLength())
		s2 = random.RandStringOfSize(b.ChannelTitle.MaxLength())
		a.ChannelTitle.SetString(s1)
		b.ChannelTitle.SetString(s2)
		if s1 != s2 && a.IsSameAs(b) {
			t.Errorf("Should not match:%s, %s", a.ChannelTitle.String(), b.ChannelTitle.String())
		}
		b.ChannelTitle.SetString(s1)

		p1 := *primitives.RandomPublicKey()
		p2 := *primitives.RandomPublicKey()
		a.LV1PublicKey = p1
		b.LV1PublicKey = p2
		if !p1.IsSameAs(&p2) && a.IsSameAs(b) {
			t.Error("Should not match")
		}
		b.LV1PublicKey = p1

		p1 = *primitives.RandomPublicKey()
		p2 = *primitives.RandomPublicKey()
		a.LV2PublicKey = p1
		b.LV2PublicKey = p2
		if !p1.IsSameAs(&p2) && a.IsSameAs(b) {
			t.Error("Should not match")
		}
		b.LV2PublicKey = p1

		p1 = *primitives.RandomPublicKey()
		p2 = *primitives.RandomPublicKey()
		a.LV3PublicKey = p1
		b.LV3PublicKey = p2
		if !p1.IsSameAs(&p2) && a.IsSameAs(b) {
			t.Error("Should not match")
		}
		b.LV3PublicKey = p1

		p1 = *primitives.RandomPublicKey()
		p2 = *primitives.RandomPublicKey()
		a.ContentSingingKey = p1
		b.ContentSingingKey = p2
		if !p1.IsSameAs(&p2) && a.IsSameAs(b) {
			t.Error("Should not match")
		}
		b.ContentSingingKey = p1

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
		a.ContentChainID = h1
		b.ContentChainID = h2
		if !h1.IsSameAs(&h2) && a.IsSameAs(b) {
			t.Error("Should not match")
		}
		b.ContentChainID = h1

		h1 = *primitives.RandomHash()
		h2 = *primitives.RandomHash()
		a.RootChainID = h1
		b.RootChainID = h2
		if !h1.IsSameAs(&h2) && a.IsSameAs(b) {
			t.Error("Should not match")
		}
		b.RootChainID = h1

		h1 = *primitives.RandomHash()
		h2 = *primitives.RandomHash()
		a.ManagementChainID = h1
		b.ManagementChainID = h2
		if !h1.IsSameAs(&h2) && a.IsSameAs(b) {
			t.Error("Should not match")
		}
		b.ManagementChainID = h1

		c1 := *RandomContentList(10)
		c2 := *RandomContentList(10)
		a.Content = c1
		b.Content = c2
		if !c1.IsSameAs(&c2) && a.IsSameAs(b) {
			t.Error("Should not match")
		}
		b.Content = c1

		pl1 := *RandomManyPlayList(10)
		pl2 := *RandomManyPlayList(10)
		a.Playlist = pl1
		b.Playlist = pl2
		if !pl1.IsSameAs(&pl2) && a.IsSameAs(b) {
			t.Error("Should not match")
		}
		b.Playlist = pl1

		hl1 := *primitives.RandomHashList(10)
		hl2 := *primitives.RandomHashList(10)
		a.SuggestedChannel = hl1
		b.SuggestedChannel = hl2
		if !hl1.IsSameAs(&hl2) && a.IsSameAs(b) {
			t.Error("Should not match")
		}
		b.SuggestedChannel = hl1

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

func TestBadUnmarshalChannel(t *testing.T) {
	badData := []byte{}

	n := NewChannel()

	_, err := n.UnmarshalBinaryData(badData)
	if err == nil {
		t.Error("Should panic or error out")
	}
}
