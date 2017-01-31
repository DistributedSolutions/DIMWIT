package primitives

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"

	//"github.com/DistributedSolutions/DIMWIT/common/primitives/random"
	ed "golang.org/x/crypto/ed25519"
)

type PrivateKey struct {
	Secret [ed.PrivateKeySize]byte
	Public PublicKey
}

func RandomPrivateKey() (*PrivateKey, error) {
	return GeneratePrivateKey()
}

// GeneratePrivateKey creates new PrivateKey / PublciKey pair or returns error
func GeneratePrivateKey() (*PrivateKey, error) {
	pk := new(PrivateKey)
	err := pk.generatePrivateKey(rand.Reader)
	return pk, err
}

func GeneratePrivateKeyFromHex(h string) (*PrivateKey, error) {
	pk := new(PrivateKey)
	data, err := hex.DecodeString(h)
	if err != nil {
		return nil, err
	}
	err = pk.SetBytes(data)
	return pk, err
}

func GeneratePrivateKeyFromBytes(data []byte) (*PrivateKey, error) {
	pk := new(PrivateKey)
	err := pk.SetBytes(data)
	return pk, err
}

func (pk *PrivateKey) generatePrivateKey(r io.Reader) error {
	pub, priv, err := ed.GenerateKey(r)
	if err != nil {
		return err
	}

	copy(pk.Secret[:ed.PrivateKeySize], priv[:ed.PrivateKeySize])
	err = pk.Public.SetBytes(pub[:])
	return err
}

func (pk *PrivateKey) SetBytes(sec []byte) error {
	if len(sec) < ed.PrivateKeySize/2 {
		return fmt.Errorf("Wrong size idiot. I need %d bytes, you gave me %d", ed.PrivateKeySize/2, len(sec))
	}

	copy(pk.Secret[:ed.PrivateKeySize/2], sec[:ed.PrivateKeySize/2])
	buf := new(bytes.Buffer)
	buf.Write(sec[:32])

	err := pk.generatePrivateKey(buf)
	return err
}

func (pk *PrivateKey) Sign(msg []byte) []byte {
	var p ed.PrivateKey
	copy(p[:ed.PrivateKeySize], pk.Secret[:ed.PrivateKeySize])
	sig := ed.Sign(p, msg)
	return sig
}

func (a *PrivateKey) IsSameAs(b *PrivateKey) bool {
	if bytes.Compare(a.Secret[:], b.Secret[:]) != 0 {
		return false
	}

	if bytes.Compare(a.Public.Bytes(), b.Public.Bytes()) != 0 {
		return false
	}

	return true
}

func (pk *PrivateKey) UnmarshalBinary(data []byte) (err error) {
	_, err = pk.UnmarshalBinaryData(data)
	return
}

func (pk *PrivateKey) UnmarshalBinaryData(data []byte) (newData []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("A panic has occurred while unmarshaling: %s", r)
			return
		}
	}()

	newData = data
	copy(pk.Secret[:ed.PrivateKeySize], newData[:ed.PrivateKeySize])

	newData = newData[ed.PrivateKeySize:]
	newData, err = pk.Public.UnmarshalBinaryData(newData)

	return
}

func (pk *PrivateKey) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	buf.Write(pk.Secret[:])

	data, err := pk.Public.MarshalBinary()
	if err != nil {
		return nil, err
	}

	buf.Write(data)
	return buf.Next(buf.Len()), nil
}
