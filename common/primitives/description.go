package primitives

import (
	"fmt"

	"github.com/DistributedSolutions/DIMWIT/common/constants"
	"github.com/DistributedSolutions/DIMWIT/common/primitives/random"
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

func (d *LongDescription) Empty() bool {
	if d.String() == "" {
		return true
	}
	return false
}

func (d *LongDescription) SetString(description string) error {
	if len(description) > d.MaxLength() {
		return fmt.Errorf("Description given is too long, length must be under %d, given length is %d",
			d.MaxLength(), len(description))
	}

	*d = LongDescription(description)
	return nil
}

func (d *LongDescription) String() string {
	return string(*d)
}

func (a *LongDescription) IsSameAs(b *LongDescription) bool {
	return a.String() == b.String()
}

func (d *LongDescription) MaxLength() int {
	return constants.LONG_DESCRIPTION_MAX_LENGTH
}

func (d *LongDescription) MarshalBinary() ([]byte, error) {
	return MarshalStringToBytes(d.String(), d.MaxLength())
}

func (d *LongDescription) UnmarshalBinary(data []byte) error {
	_, err := d.UnmarshalBinaryData(data)
	return err
}

func (d *LongDescription) UnmarshalBinaryData(data []byte) (newData []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("[LongDesc] A panic has occurred while unmarshaling: %s", r)
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

	return newData, nil
}

func RandomLongDescription() *LongDescription {
	l, _ := NewLongDescription("")
	l.SetString(random.RandStringOfSize(l.MaxLength()))

	return l
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

func (d *ShortDescription) Empty() bool {
	if d.String() == "" {
		return true
	}
	return false
}

func (d *ShortDescription) SetString(description string) error {
	if len(description) > d.MaxLength() {
		return fmt.Errorf("Description given is too long, length must be under %d, given length is %d",
			d.MaxLength(), len(description))
	}

	*d = ShortDescription(description)
	return nil
}

func (d *ShortDescription) String() string {
	return string(*d)
}

func (d *ShortDescription) MaxLength() int {
	return constants.SHORT_DESCRIPTION_MAX_LENGTH
}

func (a *ShortDescription) IsSameAs(b *ShortDescription) bool {
	return a.String() == b.String()
}

func (d *ShortDescription) MarshalBinary() ([]byte, error) {
	return MarshalStringToBytes(d.String(), d.MaxLength())
}

func (d *ShortDescription) UnmarshalBinary(data []byte) error {
	_, err := d.UnmarshalBinaryData(data)
	return err
}

func (d *ShortDescription) UnmarshalBinaryData(data []byte) (newData []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("[ShortDesc] A panic has occurred while unmarshaling: %s", r)
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

	return newData, nil
}

func RandomShortDescription() *ShortDescription {
	l, _ := NewShortDescription("")
	l.SetString(random.RandStringOfSize(l.MaxLength()))

	return l
}
