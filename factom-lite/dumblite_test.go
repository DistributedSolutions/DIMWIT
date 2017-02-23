package lite_test

import (
	. "github.com/DistributedSolutions/DIMWIT/factom-lite"
	"testing"
)

func TestGetByHeight(t *testing.T) {
	l := NewDumbLite()
	// I am real
	var _ = l

	/* _, err := l.GrabAllEntriesAtHeight(1)
	if err != nil {
		t.Error(err)
	}*/

}

func TestFake(t *testing.T) {
	l := NewFakeDumbLite()
	// I am a fake
	var _ = l
}
