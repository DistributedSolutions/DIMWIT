package primitives

import (
	"bytes"
	"fmt"

	"github.com/DistributedSolutions/DIMWIT/common/constants"
	"github.com/DistributedSolutions/DIMWIT/common/primitives/random"
)

type Image struct {
	imgType byte
	length  uint32
	image   []byte
}

func NewImage(imageBytes []byte, imageType byte) *Image {
	i := new(Image)
	i.image = imageBytes
	i.imgType = imageType
	i.length = uint32(len(imageBytes))

	return i
}

func RandomImage() *Image {
	data := random.RandByteSlice()
	i := NewImage(data, constants.IMAGE_JPEG)
	return i
}

func (i *Image) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	buf.Write([]byte{i.imgType})

	data := Uint32ToBytes(i.length)
	buf.Write(data)

	buf.Write(i.image)

	return buf.Next(buf.Len()), nil
}

func (i *Image) UnmarshalBinary(data []byte) error {
	_, err := i.UnmarshalBinaryData(data)
	return err
}

func (i *Image) UnmarshalBinaryData(data []byte) (newData []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("A panic has occurred while unmarshaling: %s", r)
			return
		}
	}()

	newData = data

	i.imgType = newData[0]
	newData = newData[1:]

	u, err := BytesToUint32(newData[:4])
	if err != nil {
		return data, err
	}
	i.length = u
	newData = newData[4:]

	i.image = newData[:i.length]
	newData = newData[i.length:]
	return
}

func (a *Image) IsSameAs(b *Image) bool {
	if a.imgType != b.imgType {
		return false
	}

	if a.length != b.length {
		return false
	}

	return BytesIsSame(a.image, b.image)
}
