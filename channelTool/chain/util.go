package chain

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"strconv"

	"github.com/DistributedSolutions/DIMWIT/common/constants"
)

type generalChainCreate interface {
	upToNonce() []byte
	getNonce() []byte
}

func FindValidNonce(i generalChainCreate) []byte {
	upToNonce := i.upToNonce()
	var count int
	count = 000
	exit := false
	for exit == false {
		count++
		exit = checkNonce(upToNonce, count)

	}

	fmt.Println(count)
	return []byte(strconv.Itoa(count))
}

func checkNonce(upToNonce []byte, nonceInt int) bool {
	buf := new(bytes.Buffer)
	buf.Write(upToNonce)

	nonce := []byte(strconv.Itoa(nonceInt))
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
