package primitives

import (
	"bytes"
	"fmt"

	"github.com/DistributedSolutions/DIMWIT/common/constants"
	"github.com/DistributedSolutions/DIMWIT/common/primitives/random"
)

type TagList struct {
	max    uint32 // Max amount of tags
	length uint32
	tags   []Tag
}

func NewTagList(max uint32) *TagList {
	tl := new(TagList)
	tl.tags = make([]Tag, 0)
	tl.max = max

	return tl
}

func RandomTagList(max uint32) *TagList {
	tl := NewTagList(max)
	l := random.RandomUInt32Between(0, max)
	tl.tags = make([]Tag, l)

	for i := range tl.tags {
		tl.tags[i] = *(RandomTag())
	}

	tl.length = l
	return tl
}

func (a *TagList) IsSameAs(b *TagList) bool {
	if a.max != b.max {
		return false
	}

	if a.length != b.length {
		return false
	}

	for i, t := range a.tags {
		if !t.IsSameAs(&(b.tags[i])) {
			return false
		}
	}

	return true
}

func (tl *TagList) AddTag(t *Tag) error {
	if uint32(len(tl.tags)) >= tl.max {
		return fmt.Errorf("Already at max tags, remove one to add another")
	}

	tl.tags = append(tl.tags, *t)
	tl.length++

	return nil
}

func (tl *TagList) GetTags() []Tag {
	return tl.tags
}

func (tl *TagList) Has(t *Tag) (int, bool) {
	for i, tt := range tl.tags {
		if t.IsSameAs(&tt) {
			return i, true
		}
	}

	return -1, false
}

func (tl *TagList) RemoveTagByName(t *Tag) error {
	i, has := tl.Has(t)
	if i == -1 || !has {
		return fmt.Errorf("Tag not found")
	}

	tl.tags = append(tl.tags[:i], tl.tags[i+1:]...)
	tl.length--
	return nil
}

func (tl *TagList) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	data, err := Uint32ToBytes(tl.max)
	if err != nil {
		return nil, err
	}

	buf.Write(data)

	data, err = Uint32ToBytes(tl.length)
	if err != nil {
		return nil, err
	}

	buf.Write(data)

	for _, t := range tl.tags {
		data, err := t.MarshalBinary()
		if err != nil {
			return nil, err
		}
		buf.Write(data)
	}

	return buf.Next(buf.Len()), nil
}

func (tl *TagList) UnmarshalBinary(data []byte) error {
	_, err := tl.UnmarshalBinaryData(data)
	return err
}

func (tl *TagList) UnmarshalBinaryData(data []byte) (newData []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("A panic has occurred while unmarshaling: %s", r)
			return
		}
	}()

	newData = data
	u, err := BytesToUint32(newData[:4])
	if err != nil {
		return data, err
	}

	newData = newData[4:]
	tl.max = u

	l, err := BytesToUint32(newData[:4])
	if err != nil {
		return data, err
	}

	newData = newData[4:]
	tl.length = l

	tl.tags = make([]Tag, tl.length)

	var i uint32 = 0
	for ; i < tl.length; i++ {
		newData, err = tl.tags[i].UnmarshalBinaryData(newData)
		if err != nil {
			return data, err
		}

	}

	return
}

type Tag string

func NewTag(tag string) (*Tag, error) {
	d := new(Tag)

	err := d.SetString(tag)
	if err != nil {
		return nil, err
	}

	return d, nil
}

func (d *Tag) SetString(tag string) error {
	if len(tag) > d.MaxLength() {
		return fmt.Errorf("Tag name given is too long, length must be under %d, given length is %d",
			d.MaxLength(), len(tag))
	}

	*d = Tag(tag)
	return nil
}

func (d *Tag) String() string {
	return string(*d) //fmt.Sprint(*d)
}

func (d *Tag) MaxLength() int {
	return constants.TAG_MAX_LENGTH
}

func (a *Tag) IsSameAs(b *Tag) bool {
	return a.String() == b.String()
}

func (d *Tag) MarshalBinary() ([]byte, error) {
	return MarshalStringToBytes(d.String(), d.MaxLength())
}

func (d *Tag) UnmarshalBinary(data []byte) error {
	_, err := d.UnmarshalBinaryData(data)
	return err
}

func (d *Tag) UnmarshalBinaryData(data []byte) (newData []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("A panic has occurred while unmarshaling: %s", r)
			return
		}
	}()

	newData = data
	str, newData, err := UnmarshalStringFromBytesData(newData, d.MaxLength())
	if err != nil {
		return data, err
	}

	err = d.SetString(str)
	if err != nil {
		return data, err
	}

	return
}

func RandomTag() *Tag {
	l, _ := NewTag("")
	index := random.RandomIntBetween(0, len(constants.ALLOWED_TAGS))
	str := constants.ALLOWED_TAGS[index]
	if len(str) > l.MaxLength() {
		str = str[:l.MaxLength()]
	}
	l.SetString(str)

	return l
}
