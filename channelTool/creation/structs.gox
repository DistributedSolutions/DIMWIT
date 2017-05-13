package creation

import (
	"github.com/FactomProject/factom"
)

type RegisterStruct struct {
	Entry *factom.Entry
}

type CreateStruct struct {
	Chain    *factom.Chain
	ExtIDs   [][]byte
	endExtID int
}

func (c CreateStruct) upToNonce() []byte {
	return upToNonce(c.ExtIDs)
}
