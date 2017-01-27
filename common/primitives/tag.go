package primitives

import (
	"fmt"

	"github.com/DistributedSolutions/DIMWIT/common/constants"
)

type TagList struct {
	Length int
	Tags   []Tag
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
	if len(tag) > constants.TAG_MAX_LENGTH {
		return fmt.Errorf("Description given is too long, length must be under %d, given length is %d",
			constants.TAG_MAX_LENGTH, len(tag))
	}

	*d = Tag(tag)
	return nil
}

func (d *Tag) String() string {
	return string(*d)
}
