package primitives

import (
	"encoding/hex"
	"fmt"

	"github.com/DistributedSolutions/DIMWIT/common/constants"
	"github.com/DistributedSolutions/DIMWIT/common/primitives/random"
)

type HashList struct {
	Length int
	List   []Hash
}

type Hash [constants.HASH_LENGTH]byte

func BytesToHash(b []byte) (*Hash, error) {
	h := new(Hash)
	err := h.SetBytes(b)
	if err != nil {
		return nil, err
	}

	return h, nil
}

func HexToHash(he string) (*Hash, error) {
	data, err := hex.DecodeString(he)
	if err != nil {
		return nil, err
	}

	h := new(Hash)
	err = h.SetBytes(data)
	if err != nil {
		return nil, err
	}

	return h, nil
}

func (h *Hash) Bytes() []byte {
	ni := make([]byte, h.Length())
	copy(ni, h[:])
	return ni
}

func (h *Hash) SetBytes(ni []byte) error {
	if len(ni) != h.Length() {
		return fmt.Errorf("Length is invalid, must be of length %d", h.Length())
	}

	copy(h[:], ni)
	return nil
}

func (h *Hash) String() string {
	return hex.EncodeToString(h.Bytes())
}

func RandomHash() *Hash {
	h := new(Hash)
	h.SetBytes(random.RandByteSliceOfSize(h.Length()))
	return h
}

func (h *Hash) Length() int {
	return constants.HASH_LENGTH
}

func (h *Hash) MarshalBinary() ([]byte, error) {
	return h.Bytes(), nil
}

func (h *Hash) UnmarshalBinary(data []byte) error {
	_, err := h.UnmarshalBinaryData(data)
	return err
}

func (h *Hash) UnmarshalBinaryData(data []byte) (newData []byte, err error) {
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

func (a *Hash) IsSameAs(b *Hash) bool {
	adata := a.Bytes()
	bdata := b.Bytes()
	for i := 0; i < a.Length(); i++ {
		if adata[i] != bdata[i] {
			return false
		}
	}

	return true
}
