package primitives

import (
	"fmt"

	"github.com/DistributedSolutions/DIMWIT/common/constants"
)

type LongDescription string

func NewLongDescription(description string) (*LongDescription, error) {
	d := new(LongDescription)

	err := d.SetString(description)
	if err != nil {
		return nil, err
	}

	return d, nil
}

func (d *LongDescription) SetString(description string) error {
	if len(description) > constants.LONG_DESCRIPTION_MAX_LENGTH {
		return fmt.Errorf("Description given is too long, length must be under %d, given length is %d",
			constants.LONG_DESCRIPTION_MAX_LENGTH, len(description))
	}

	*d = LongDescription(description)
	return nil
}

func (d *LongDescription) String() string {
	return string(*d)
}

type ShortDescription string

func NewShortDescription(description string) (*ShortDescription, error) {
	s := new(ShortDescription)

	err := s.SetString(description)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (d *ShortDescription) SetString(description string) error {
	if len(description) > constants.SHORT_DESCRIPTION_MAX_LENGTH {
		return fmt.Errorf("Description given is too long, length must be under %d, given length is %d",
			constants.SHORT_DESCRIPTION_MAX_LENGTH, len(description))
	}

	*d = ShortDescription(description)
	return nil
}

func (d *ShortDescription) String() string {
	return string(*d)
}
