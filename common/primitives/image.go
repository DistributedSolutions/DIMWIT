package primitives

import (
	"bytes"
	//"encoding/hex"
	"encoding/json"
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
	i.length = uint32(len(i.image))
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
	return i.length
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

	return bytes.Compare(a.image, b.image) == 0
}

func (i *Image) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		ImgType byte   `json:"imgtype"`
		Length  uint32 `json:"length"`
		Image   string `json:"image"`
	}{
		ImgType: i.imgType,
		Length:  i.length,
		Image:   string(i.image),
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
	/*img, err := hex.DecodeString(obj.Image)
	if err != nil {
		return err
	}*/
	img := []byte(obj.Image)
	i.SetImage(img)
	return nil
}
