package primitives

import (
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/DistributedSolutions/DIMWIT/common/constants"
	"github.com/DistributedSolutions/DIMWIT/common/primitives/random"
)

type InfoHash [constants.INFOHASH_BYTES_LENGTH]byte

func BytesToInfoHash(b []byte) (*InfoHash, error) {
	i := new(InfoHash)
	err := i.SetBytes(b)
	if err != nil {
		return nil, err
	}

	return i, nil
}

func HexToInfoHash(h string) (*InfoHash, error) {
	data, err := hex.DecodeString(h)
	if err != nil {
		return nil, err
	}

	i := new(InfoHash)
	err = i.SetBytes(data)
	if err != nil {
		return nil, err
	}

	return i, nil
}

func (i *InfoHash) Empty() bool {
	for _, b := range i.Bytes() {
		if b != 0x00 {
			return false
		}
	}
	return true
}

func (i *InfoHash) Bytes() []byte {
	ni := make([]byte, i.Length())
	copy(ni, i[:])
	return ni
}

func (i *InfoHash) SetBytes(ni []byte) error {
	if len(ni) != i.Length() {
		return fmt.Errorf("Length is invalid, must be of length %d", i.Length())
	}

	copy(i[:], ni)
	return nil
}

func (i *InfoHash) String() string {
	return hex.EncodeToString(i.Bytes())
}

func RandomInfoHash() *InfoHash {
	h := new(InfoHash)
	h.SetBytes(random.RandByteSliceOfSize(h.Length()))
	return h
}

func (h *InfoHash) Length() int {
	return constants.INFOHASH_BYTES_LENGTH
}

func (h *InfoHash) MarshalBinary() ([]byte, error) {
	return h.Bytes(), nil
}

func (h *InfoHash) UnmarshalBinary(data []byte) error {
	_, err := h.UnmarshalBinaryData(data)
	return err
}

func (h *InfoHash) UnmarshalBinaryData(data []byte) (newData []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("[InfoHash] A panic has occurred while unmarshaling: %s", r)
			return
		}
	}()

	newData = data
	if len(newData) < h.Length() {
		err = fmt.Errorf("Length is invalid, must be of length %d, found length %d", h.Length(), len(newData))
		return
	}

	err = h.SetBytes(newData[:h.Length()])
	if err != nil {
		return
	}
	newData = newData[h.Length():]
	return
}

func (a *InfoHash) IsSameAs(b *InfoHash) bool {
	adata := a.Bytes()
	bdata := b.Bytes()
	for i := 0; i < a.Length(); i++ {
		if adata[i] != bdata[i] {
			return false
		}
	}

	return true
}

func (h *InfoHash) MarshalJSON() ([]byte, error) {
	return json.Marshal(h.String())
}

func (h *InfoHash) UnmarshalJSON(b []byte) error {
	var hexS string
	if err := json.Unmarshal(b, &hexS); err != nil {
		return err
	}
	data, err := hex.DecodeString(hexS)
	if err != nil {
		return err
	}
	h.SetBytes(data)
	return nil
}
