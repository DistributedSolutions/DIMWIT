package primitives_test

import (
	"encoding/json"
	"fmt"
	"testing"

	. "github.com/DistributedSolutions/DIMWIT/common/primitives"
	//"github.com/DistributedSolutions/DIMWIT/common/primitives/random"
)

var _ = fmt.Sprintf("")

func TestImage(t *testing.T) {
	for i := 0; i < 1000; i++ {
		h := RandomImage()
		if i%100 == 0 {
			h = RandomHugeImage()
		}
		data, _ := h.MarshalBinary()

		n := new(Image)
		newData, err := n.UnmarshalBinaryData(data)
		if err != nil {
			t.Error(err)
		}

		if !h.IsSameAs(n) {
			t.Error("Failed, should be same")
		}

		if len(newData) != 0 {
			t.Error("Failed, should have no bytes left")
		}

		if h.Empty() && h.GetImageSize() != 0 {
			t.Error("Should not be empty")
		}

		if i > 10 {
			continue
		}

		jdata, err := json.Marshal(h)
		if err != nil {
			t.Error(err)
		}

		j := new(Image)
		err = json.Unmarshal(jdata, j)
		if err != nil {
			t.Error(err)
		}

		if !h.IsSameAs(j) {
			t.Errorf("Should be same. ")
		}
	}
}

func TestDiffImage(t *testing.T) {
	same := 0
	for i := 0; i < 1000; i++ {
		a := RandomImage()
		b := RandomImage()
		if a.IsSameAs(b) {
			same++
		}
	}
	if same > 15 {
		t.Error("More than 15 are the same, it is totally random, so it is likely the IsSameAs is broken.")
	}

	a := RandomImage()
	b := RandomImage()

	if len(a.GetImage()) < 1 {
		a.SetImage([]byte("Random"))
	}

	b.SetImage(a.GetImage()[1:]) // Failed here
	if a.IsSameAs(b) {
		t.Error("Not same")
	}

	a.SetImage([]byte("A test"))
	b.SetImage([]byte("A test2"))
	if a.IsSameAs(b) {
		t.Error("Not same")
	}

	b.SetImage(a.GetImage())
	a.SetImageType(0x02)
	b.SetImageType(0x00)
	if a.IsSameAs(b) {
		t.Error("Not same")
	}

	a.GetImageType()
}

func TestEmptyImage(t *testing.T) {
	i := new(Image)
	if !i.Empty() {
		t.Error("Should be empty")
	}
}

func TestBadUnmarshalImage(t *testing.T) {
	badData := []byte{}

	n := new(Image)

	_, err := n.UnmarshalBinaryData(badData)
	if err == nil {
		t.Error("Should panic or error out")
	}

	badJson := ""
	err = n.UnmarshalJSON([]byte(badJson))
	if err == nil {
		t.Error("Should panic or error out")
	}

	badJson = `{"imgtype":2,"length":6,"image":"T12074657374"}`
	err = n.UnmarshalJSON([]byte(badJson))
	if err == nil {
		t.Error("Should panic or error out")
	}
}
