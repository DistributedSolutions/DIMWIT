package creation

import ()

type RegisterStruct struct {
	Entry *factom.Entry
}

type CreateStruct struct {
	Chain  *factom.Chain
	ExtIDs [][]byte
}

func (c CreateStruct) upToNonce() []byte {
	return upToNonce(c.ExtIDs, c.endExtID)
}
