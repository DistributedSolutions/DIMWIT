package primitives

import (
	"encoding/hex"
	"fmt"

	"github.com/DistributedSolutions/DIMWIT/common/constants"
)

type InfoHash [constants.INFOHASH_LENGTH]byte

func HexToHash(h string) (*InfoHash, error) {
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

func (i *InfoHash) Bytes() []byte {
	ni := make([]byte, constants.INFOHASH_LENGTH)
	copy(ni, i[:])
	return ni
}

func (i *InfoHash) SetBytes(ni []byte) error {
	if len(ni) != constants.INFOHASH_LENGTH {
		return fmt.Errorf("Length is invalid, must be of length %d", constants.INFOHASH_LENGTH)
	}

	copy(i[:], ni)
	return nil
}
