package primitives

import (
	"bytes"
	"fmt"

	"github.com/DistributedSolutions/DIMWIT/common/constants"
	"github.com/DistributedSolutions/DIMWIT/common/primitives/random"
)

type TagList struct {
	Max  uint32 `json:"max"` // Max amount of tags
	Tags []Tag  `json:"tags"`
}

func NewTagList(max uint32) *TagList {
	tl := new(TagList)
	tl.Tags = make([]Tag, 0)
	tl.Max = max

	return tl
}

func RandomTagList(max uint32) *TagList {
	tl := NewTagList(max)
	l := uint32(1)
	tl.Tags = make([]Tag, l)

	c := uint32(0)
	for i := uint32(0); i < l; i++ {
		tempTag := RandomTag()
		_, b := tl.HasTag(tempTag)
		if !b {
			tl.Tags[c] = *tempTag
			c++
		}
	}
	tl.Tags = tl.Tags[:c]
	return tl
}

func (a *TagList) Combine(b *TagList) *TagList {
	t := NewTagList(a.Max)
	tl := make([]Tag, 0)
	if uint32(len(a.Tags))+uint32(len(b.Tags)) < a.Max {
		tl = append(a.Tags, b.Tags...)
	} else if uint32(len(b.Tags)) >= a.Max {
		tl = b.Tags[:a.Max]
	} else {
		amt := a.Max - uint32(len(b.Tags))
		tl = append(a.Tags[:amt], b.Tags...)
	}

	t.Tags = tl
	return t
}

func (d *TagList) MaxTags() uint32 {
	return d.Max
}

func (d *TagList) Empty() bool {
	if len(d.Tags) == 0 {
		return true
	}
	return false
}

func (a *TagList) IsSameAs(b *TagList) bool {
	if a.Max != b.Max {
		return false
	}

	if len(a.Tags) != len(b.Tags) {
		return false
	}

	for i, t := range a.Tags {
		if !t.IsSameAs(&(b.Tags[i])) {
			return false
		}
	}

	return true
}

func (tl *TagList) AddTagByName(t string) error {
	tag, err := NewTag(t)
	if err != nil {
		return err
	}
	return tl.AddTag(tag)
}

func (tl *TagList) AddTag(t *Tag) error {
	if uint32(len(tl.Tags)) >= tl.Max {
		return fmt.Errorf("Already at max tags, remove one to add another")
	}

	tl.Tags = append(tl.Tags, *t)

	return nil
}

func (tl *TagList) GetTags() []Tag {
	return tl.Tags
}

func (tl *TagList) GetTagsAsStringArr() []string {
	arr := make([]string, len(tl.Tags))
	for i := 0; i < len(tl.Tags); i++ {
		arr[i] = tl.Tags[i].String()
	}
	return arr
}

func (tl *TagList) Has(t string) (int, bool) {
	tag, err := NewTag(t)
	if err != nil {
		return -1, false
	}
	return tl.HasTag(tag)
}

func (tl *TagList) HasTag(t *Tag) (int, bool) {
	for i, tt := range tl.Tags {
		if t.IsSameAs(&tt) {
			return i, true
		}
	}

	return -1, false
}

func (tl *TagList) SetTagTo(index int, tag string) error {
	if len(tl.Tags) <= index {
		return fmt.Errorf("Tag not found")
	}

	tl.Tags[index].SetString(tag)
	return nil
}

func (tl *TagList) RemoveTagByName(t string) error {
	tag, err := NewTag(t)
	if err != nil {
		return err
	}
	return tl.RemoveTag(tag)
}

func (tl *TagList) RemoveTag(t *Tag) error {
	i, has := tl.HasTag(t)
	if i == -1 || !has {
		return fmt.Errorf("Tag not found")
	}

	tl.Tags = append(tl.Tags[:i], tl.Tags[i+1:]...)
	return nil
}

func (tl *TagList) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	data := Uint32ToBytes(tl.Max)

	buf.Write(data)

	data = Uint32ToBytes(uint32(len(tl.Tags)))

	buf.Write(data)

	for _, t := range tl.Tags {
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
			err = fmt.Errorf("[TagList] A panic has occurred while unmarshaling: %s", r)
			return
		}
	}()

	newData = data
	u, err := BytesToUint32(newData[:4])
	if err != nil {
		return data, err
	}

	newData = newData[4:]
	tl.Max = u

	l, err := BytesToUint32(newData[:4])
	if err != nil {
		return data, err
	}

	newData = newData[4:]
	tl.Tags = make([]Tag, l)

	var i uint32 = 0
	for ; i < l; i++ {
		t := new(Tag)
		newData, err = t.UnmarshalBinaryData(newData)
		if err != nil {
			return data, err
		}

		tl.Tags[i] = *t
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

func (d *Tag) Empty() bool {
	if d.String() == "" {
		return true
	}
	return false
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
			err = fmt.Errorf("[Tag] A panic has occurred while unmarshaling: %s", r)
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
