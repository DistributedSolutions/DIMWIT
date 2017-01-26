package primitives

import (
	"encoding/hex"
	"fmt"

	"github.com/DistributedSolutions/DIMWIT/common/constants"
)

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
	ni := make([]byte, constants.HASH_LENGTH)
	copy(ni, h[:])
	return ni
}

func (h *Hash) SetBytes(ni []byte) error {
	if len(ni) != constants.HASH_LENGTH {
		return fmt.Errorf("Length is invalid, must be of length %d", constants.HASH_LENGTH)
	}

	copy(h[:], ni)
	return nil
}

func (h *Hash) String() string {
	return hex.EncodeToString(h.Bytes())
}
