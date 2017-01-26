package primitives

import ()

type Image struct {
	Length int
	Type   byte
	Image  []byte
}

func NewImage(imageBytes []byte, imageType byte) *Image {
	i := new(Image)
	i.Image = imageBytes
	i.Type = imageType
	i.Length = len(imageBytes)

	return i
}
