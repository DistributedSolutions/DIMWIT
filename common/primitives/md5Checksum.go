package primitives

import (
	"encoding/hex"
	"fmt"

	"github.com/DistributedSolutions/DIMWIT/common/constants"
	"github.com/DistributedSolutions/DIMWIT/common/primitives/random"
)

type MD5Checksum [constants.MD5_CHECKSUM_LENGTH]byte

func BytesToMD5Checksum(b []byte) (*MD5Checksum, error) {
	m := new(MD5Checksum)
	err := m.SetBytes(b)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func HexToMD5Checksum(he string) (*MD5Checksum, error) {
	data, err := hex.DecodeString(he)
	if err != nil {
		return nil, err
	}

	m := new(MD5Checksum)
	err = m.SetBytes(data)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (m *MD5Checksum) Bytes() []byte {
	ni := make([]byte, m.Length())
	copy(ni, m[:])
	return ni
}

func (m *MD5Checksum) SetBytes(ni []byte) error {
	if len(ni) != m.Length() {
		return fmt.Errorf("Length is invalid, must be of length %d", m.Length())
	}

	copy(m[:], ni)
	return nil
}

func (m *MD5Checksum) String() string {
	return hex.EncodeToString(m.Bytes())
}

func RandomMD5() *MD5Checksum {
	h := new(MD5Checksum)
	h.SetBytes(random.RandByteSliceOfSize(h.Length()))
	return h
}

func (h *MD5Checksum) Length() int {
	return constants.MD5_CHECKSUM_LENGTH
}

func (h *MD5Checksum) MarshalBinary() ([]byte, error) {
	return h.Bytes(), nil
}

func (h *MD5Checksum) UnmarshalBinary(data []byte) error {
	_, err := h.UnmarshalBinaryData(data)
	return err
}

func (h *MD5Checksum) UnmarshalBinaryData(data []byte) (newData []byte, err error) {
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

func (a *MD5Checksum) IsSameAs(b *MD5Checksum) bool {
	adata := a.Bytes()
	bdata := b.Bytes()
	for i := 0; i < a.Length(); i++ {
		if adata[i] != bdata[i] {
			return false
		}
	}

	return true
}
