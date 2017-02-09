package primitives

import (
	"encoding/hex"
	"fmt"

	"github.com/DistributedSolutions/DIMWIT/common/primitives/random"
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

func (p *PublicKey) Empty() bool {
	for _, b := range p.Bytes() {
		if b != 0x00 {
			return false
		}
	}
	return true
}

func (p *PublicKey) Bytes() []byte {
	b := make([]byte, p.Length())
	copy(b, p[:])
	return b
}

func (p *PublicKey) FixedBytes() [ed.PublicKeySize]byte {
	var b [ed.PublicKeySize]byte
	copy(b[:], p[:])
	return b
}

func (p *PublicKey) SetBytes(b []byte) error {
	if len(b) != ed.PublicKeySize {
		return fmt.Errorf("Length is invalid, must be of length %d", p.Length())
	}

	copy(p[:], b)
	return nil
}

func (p *PublicKey) String() string {
	return hex.EncodeToString(p.Bytes())
}

func RandomPublicKey() *PublicKey {
	h := new(PublicKey)
	h.SetBytes(random.RandByteSliceOfSize(h.Length()))
	return h
}

func (h *PublicKey) Length() int {
	return ed.PublicKeySize
}

func (h *PublicKey) MarshalBinary() ([]byte, error) {
	return h.Bytes(), nil
}

func (h *PublicKey) UnmarshalBinary(data []byte) error {
	_, err := h.UnmarshalBinaryData(data)
	return err
}

func (h *PublicKey) UnmarshalBinaryData(data []byte) (newData []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("[PubKey] A panic has occurred while unmarshaling: %s", r)
			return
		}
	}()

	newData = data
	if len(newData) < h.Length() {
		err = fmt.Errorf("Length is invalid, must be of length %d, found length %d", h.Length(), len(newData))
		return
	}

	err = h.SetBytes(newData[:h.Length()])
	if err != nil {
		return
	}
	newData = newData[h.Length():]
	return
}

func (a *PublicKey) IsSameAs(b *PublicKey) bool {
	adata := a.Bytes()
	bdata := b.Bytes()
	for i := 0; i < a.Length(); i++ {
		if adata[i] != bdata[i] {
			return false
		}
	}

	return true
}

func (p *PublicKey) Verify(msg []byte, sig []byte) bool {
	return ed.Verify(p.Bytes(), msg, sig)
}
