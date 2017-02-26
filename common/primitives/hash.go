package primitives

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/DistributedSolutions/DIMWIT/common/constants"
	"github.com/DistributedSolutions/DIMWIT/common/primitives/random"
)

type HashList struct {
	List []Hash `json:"hashlist"`
}

type Hash [constants.HASH_BYTES_LENGTH]byte

func RandomHashList(max uint32) *HashList {
	h := NewHashList()
	l := random.RandomUInt32Between(0, max)
	h.List = make([]Hash, l)

	for i := range h.List {
		h.List[i] = *RandomHash()
	}

	return h
}

func NewZeroHash() *Hash {
	h, _ := HexToHash("0000000000000000000000000000000000000000000000000000000000000000")
	return h
}

func NewHashList() *HashList {
	h := new(HashList)
	h.List = make([]Hash, 0)

	return h
}

func (a *HashList) Combine(b *HashList) *HashList {
	x := new(HashList)
	x.List = append(a.List, b.List...)
	return x
}

func (a *HashList) Empty() bool {
	if len(a.List) == 0 {
		return true
	}
	return false
}

func (a *HashList) IsSameAs(b *HashList) bool {
	if len(a.List) != len(b.List) {
		return false
	}

	for i := range a.List {
		if a.List[i] != b.List[i] {
			return false
		}
	}

	return true
}

func (h *HashList) GetHashes() []Hash {
	return h.List
}

func (h *HashList) AddHash(hash *Hash) {
	h.List = append(h.List, *hash)
}

func (h *HashList) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	data := Uint32ToBytes(uint32(len(h.List)))

	buf.Write(data)

	for i := range h.List {
		// This cannot actually error out
		data, _ := h.List[i].MarshalBinary()
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

	h.List = make([]Hash, u)
	var i uint32
	for i = 0; i < u; i++ {
		newData, err = h.List[i].UnmarshalBinaryData(newData)
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

func (h *Hash) MarshalJSON() ([]byte, error) {
	return json.Marshal(h.String())
}

func (h *Hash) UnmarshalJSON(b []byte) error {
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
