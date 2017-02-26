package primitives

import (
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/DistributedSolutions/DIMWIT/common/constants"
	"github.com/DistributedSolutions/DIMWIT/common/primitives/random"
)

type MD5Checksum [constants.MD5_CHECKSUM_BYTES_LENGTH]byte

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

func (m *MD5Checksum) Empty() bool {
	for _, b := range m.Bytes() {
		if b != 0x00 {
			return false
		}
	}
	return true
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
	return constants.MD5_CHECKSUM_BYTES_LENGTH
}

func (h *MD5Checksum) MarshalBinary() ([]byte, error) {
	return h.Bytes(), nil
}

func (h *MD5Checksum) UnmarshalBinary(data []byte) error {
	_, err := h.UnmarshalBinaryData(data)
	return err
}

func (h *MD5Checksum) UnmarshalBinaryData(data []byte) (newData []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("[Md5] A panic has occurred while unmarshaling: %s", r)
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

func (h *MD5Checksum) MarshalJSON() ([]byte, error) {
	return json.Marshal(h.String())
}

func (h *MD5Checksum) UnmarshalJSON(b []byte) error {
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
