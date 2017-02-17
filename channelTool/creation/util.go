package creation

import (
	"bytes"
	//"crypto/rand"
	"crypto/sha256"

	"github.com/DistributedSolutions/DIMWIT/common/constants"
	"github.com/DistributedSolutions/DIMWIT/common/primitives"
)

type generalChainCreate interface {
	upToNonce() []byte
}

func FindValidNonce(i generalChainCreate) []byte {
	upToNonce := i.upToNonce()
	var count uint64
	count = 0
	exit := false
	for exit == false {
		count++
		exit = checkNonce(upToNonce, count)

	}

	data, _ := primitives.Uint64ToBytes(count)
	return data
}

func checkNonce(upToNonce []byte, nonceInt uint64) bool {
	buf := new(bytes.Buffer)
	buf.Write(upToNonce)

	nonce, _ := primitives.Uint64ToBytes(nonceInt)
	//nonce := []byte(strconv.Itoa(nonceInt))
	result := sha256.Sum256(nonce)
	buf.Write(result[:])

	result = sha256.Sum256(buf.Bytes())

	chainFront := result[:constants.CHAIN_PREFIX_LENGTH_CHECK]

	if bytes.Compare(chainFront[:constants.CHAIN_PREFIX_LENGTH_CHECK],
		constants.CHAIN_PREFIX[:constants.CHAIN_PREFIX_LENGTH_CHECK]) == 0 {
		return true
	}
	return false
}

// upToNonce is exclusive
func upToNonce(extIDs [][]byte) []byte {
	buf := new(bytes.Buffer)
	for _, e := range extIDs {
		result := sha256.Sum256(e)
		buf.Write(result[:])
	}

	return buf.Next(buf.Len())
}

func ExIDLength(exid [][]byte) int {
	length := 2
	for _, e := range exid {
		length += 2
		length += len(e)
	}

	return length
}
