package elements

import (
	"bytes"
	"time"
	//"crypto/rand"
	"crypto/sha256"

	"github.com/DistributedSolutions/DIMWIT/common/constants"
	"github.com/DistributedSolutions/DIMWIT/common/primitives"
)

func TimeStampBytes() []byte {
	ts := time.Now()
	tsData, err := ts.MarshalBinary()
	if err != nil {
		// Uh.... really?
		// return err
	}
	return tsData
}

func VersionAndType(e IElementPiece) [][]byte {
	extID := make([][]byte, 0)
	extID = append(extID, GetVersionBytes())
	extID = append(extID, e.Type())
	return extID
}

func GetVersionBytes() []byte {
	return []byte{constants.FACTOM_VERSION}
}

func GetVersion() byte {
	return constants.FACTOM_VERSION
}

// GetContentSignature leaves a signature in the body of some entries
// to signal this tool was used
func GetContentSignature() []byte {
	return []byte("Created with EZ-Tool")
}

func FindValidNonce(exceptNocne [][]byte) (nonce []byte, res []byte) {
	upToNonce := upToNonce(exceptNocne)
	var count uint64
	count = 0
	exit := false
	for exit == false {
		count++
		exit, res = checkNonce(upToNonce, count)

	}

	data, _ := primitives.Uint64ToBytes(count)
	return data, res
}

func checkNonce(upToNonce []byte, nonceInt uint64) (bool, []byte) {
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
		return true, result[:]
	}
	return false, nil
}

func upToSig(extIDs [][]byte) []byte {
	buf := new(bytes.Buffer)
	for _, e := range extIDs {
		buf.Write(e)
	}

	return buf.Next(buf.Len())
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

func howManyEntries(headerLength int, contentLength int, contentHeaderLength int) int {
	contentLength -= (constants.ENTRY_MAX_SIZE - headerLength)
	bytesPerEntry := constants.ENTRY_MAX_SIZE - contentHeaderLength
	count := 0
	for contentLength > 0 {
		contentLength -= bytesPerEntry
		count++
	}

	return count
}
