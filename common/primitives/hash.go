package primitives

import (
	"bytes"
	"encoding/hex"
	"fmt"

	"github.com/DistributedSolutions/DIMWIT/common/constants"
	"github.com/DistributedSolutions/DIMWIT/common/primitives/random"
)

type HashList struct {
	list []Hash
}

type Hash [constants.HASH_BYTES_LENGTH]byte

func RandomHashList(max uint32) *HashList {
	h := NewHashList()
	l := random.RandomUInt32Between(0, max)
	h.list = make([]Hash, l)

	for i := range h.list {
		h.list[i] = *RandomHash()
	}

	return h
}

func NewZeroHash() *Hash {
	h, _ := HexToHash("0000000000000000000000000000000000000000000000000000000000000000")
	return h
}

func NewHashList() *HashList {
	h := new(HashList)
	h.list = make([]Hash, 0)

	return h
}

func (a *HashList) Combine(b *HashList) *HashList {
	x := new(HashList)
	x.list = append(a.list, b.list...)
	return x
}

func (a *HashList) Empty() bool {
	if len(a.list) == 0 {
		return true
	}
	return false
}

func (a *HashList) IsSameAs(b *HashList) bool {
	if len(a.list) != len(b.list) {
		return false
	}

	for i := range a.list {
		if a.list[i] != b.list[i] {
			return false
		}
	}

	return true
}

func (h *HashList) GetHashes() []Hash {
	return h.list
}

func (h *HashList) AddHash(hash *Hash) {
	h.list = append(h.list, *hash)
}

func (h *HashList) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	data := Uint32ToBytes(uint32(len(h.list)))

	buf.Write(data)

	for i := range h.list {
		// This cannot actually error out
		data, _ := h.list[i].MarshalBinary()
		buf.Write(data)
	}

	return buf.Next(buf.Len()), nil
}

func (h *HashList) UnmarshalBinary(data []byte) error {
	_, err := h.UnmarshalBinaryData(data)
	return err
}

func (h *HashList) UnmarshalBinaryData(data []byte) (newData []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("[HashList] A panic has occurred while unmarshaling: %s", r)
			return
		}
	}()

	newData = data
	u, err := BytesToUint32(newData[:4])
	if err != nil {
		return data, err
	}
	newData = newData[4:]

	h.list = make([]Hash, u)
	var i uint32
	for i = 0; i < u; i++ {
		newData, err = h.list[i].UnmarshalBinaryData(newData)
		if err != nil {
			return data, err
		}
	}

	return
}

// TODO: Hashlist

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

func (a *Hash) Empty() bool {
	if NewZeroHash().IsSameAs(a) {
		return true
	}
	return false
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
	return constants.HASH_BYTES_LENGTH
}

func (h *Hash) MarshalBinary() ([]byte, error) {
	return h.Bytes(), nil
}

func (h *Hash) UnmarshalBinary(data []byte) error {
	_, err := h.UnmarshalBinaryData(data)
	return err
}

func (h *Hash) UnmarshalBinaryData(data []byte) (newData []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("[Hash] A panic has occurred while unmarshaling: %s", r)
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
