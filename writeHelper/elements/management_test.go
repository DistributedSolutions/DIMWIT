package elements_test

import (
	"testing"

	. "github.com/DistributedSolutions/DIMWIT/writeHelper/elements"
)

func TestStripEmpty(t *testing.T) {
	m := NewManageChainMetaData()
	m.StripEmpty()

	if m.Website != nil {
		t.Fail()
	}

	if m.LongDescription != nil {
		t.Fail()
	}

	if m.ShortDescription != nil {
		t.Fail()
	}

	if m.Playlist != nil {
		t.Fail() //ch.Playlist = meta.Playlist
	}

	if m.Thumbnail != nil {
		t.Fail()
	}

	if m.Banner != nil {
		t.Fail()
	}

	if m.ChannelTags != nil {
		t.Fail()
	}

	if m.SuggestedChannels != nil {
		t.Fail()
	}
}
