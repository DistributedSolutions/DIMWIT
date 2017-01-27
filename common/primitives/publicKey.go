package primitives

import (
	"encoding/hex"
	"fmt"

	ed "golang.org/x/crypto/ed25519"
)

type PublicKey [ed.PublicKeySize]byte

func PublicKeyFromBytes(b []byte) (*PublicKey, error) {
	p := new(PublicKey)
	err := p.SetBytes(b)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func PublicKeyFromHex(he string) (*PublicKey, error) {
	data, err := hex.DecodeString(he)
	if err != nil {
		return nil, err
	}

	p := new(PublicKey)
	err = p.SetBytes(data)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (p *PublicKey) Bytes() []byte {
	b := make([]byte, ed.PublicKeySize)
	copy(b, p[:])
	return b
}

func (p *PublicKey) SetBytes(b []byte) error {
	if len(b) != ed.PublicKeySize {
		return fmt.Errorf("Length is invalid, must be of length %d", ed.PublicKeySize)
	}

	copy(p[:], b)
	return nil
}

func (p *PublicKey) String() string {
	return hex.EncodeToString(p.Bytes())
}
