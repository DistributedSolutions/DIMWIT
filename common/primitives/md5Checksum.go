package primitives

import (
	"encoding/hex"
	"fmt"

	"github.com/DistributedSolutions/DIMWIT/common/constants"
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
	ni := make([]byte, constants.MD5_CHECKSUM_LENGTH)
	copy(ni, m[:])
	return ni
}

func (m *MD5Checksum) SetBytes(ni []byte) error {
	if len(ni) != constants.MD5_CHECKSUM_LENGTH {
		return fmt.Errorf("Length is invalid, must be of length %d", constants.MD5_CHECKSUM_LENGTH)
	}

	copy(m[:], ni)
	return nil
}

func (m *MD5Checksum) String() string {
	return hex.EncodeToString(m.Bytes())
}
