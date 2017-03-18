package primitives

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/DistributedSolutions/DIMWIT/common/constants"
	"github.com/DistributedSolutions/DIMWIT/common/primitives/random"
)

type Image struct {
	imgType byte
	image   []byte
}

func NewImage(imageBytes []byte, imageType byte) *Image {
	i := new(Image)
	i.image = imageBytes
	i.imgType = imageType

	return i
}

func RandomImage() *Image {
	data := random.RandByteSlice()
	i := NewImage(data, constants.IMAGE_JPEG)
	return i
}

func RandomHugeImage() *Image {
	data := random.RandByteSliceOfSize(random.RandomIntBetween(100000, 300000))
	i := NewImage(data, constants.IMAGE_JPEG)
	return i
}

func (i *Image) Empty() bool {
	if len(i.image) == 0 {
		return true
	}
	return false
}

func (i *Image) SetImage(data []byte) {
	i.image = data
}

func (i *Image) SetImageType(t byte) {
	i.imgType = t
}

func (i *Image) GetImageType() byte {
	return i.imgType
}

func (i *Image) GetImage() []byte {
	return i.image
}

func (i *Image) GetImageSize() uint32 {
	return uint32(len(i.image))
}

func (i *Image) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	buf.Write([]byte{i.imgType})

	data := Uint32ToBytes(uint32(len(i.image)))
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
			err = fmt.Errorf("[Image] A panic has occurred while unmarshaling: %s", r)
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
	newData = newData[4:]

	i.image = newData[:u]
	newData = newData[u:]
	return
}

func (a *Image) IsSameAs(b *Image) bool {
	if a.imgType != b.imgType {
		return false
	}

	if len(a.image) != len(b.image) {
		return false
	}

	return bytes.Compare(a.image, b.image) == 0
}

func (i *Image) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		ImgType byte   `json:"imgtype"`
		Length  uint32 `json:"length"`
		Image   string `json:"image"`
	}{
		ImgType: i.imgType,
		Length:  uint32(len(i.image)),
		Image:   hex.EncodeToString(i.image),
	})
}

func (i *Image) UnmarshalJSON(b []byte) error {
	obj := new(struct {
		ImgType byte   `json:"imgtype"`
		Length  uint32 `json:"length"`
		Image   string `json:"image"`
	})

	if err := json.Unmarshal(b, obj); err != nil {
		return err
	}

	i.SetImageType(obj.ImgType)
	img, err := hex.DecodeString(obj.Image)
	if err != nil {
		return err
	}
	// img := []byte(obj.Image)
	i.SetImage(img)
	return nil
}
