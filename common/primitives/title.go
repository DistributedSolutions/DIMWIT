package primitives

import (
	"fmt"

	"github.com/DistributedSolutions/DIMWIT/common/constants"
)

type Title string

func NewTitle(title string) (*Title, error) {
	d := new(Title)

	err := d.SetString(title)
	if err != nil {
		return nil, err
	}

	return d, nil
}

func (d *Title) SetString(title string) error {
	if len(title) > d.MaxLength() {
		return fmt.Errorf("Description given is too long, length must be under %d, given length is %d",
			d.MaxLength(), len(title))
	}

	*d = Title(title)
	return nil
}

func (d *Title) String() string {
	return string(*d)
}

func (d *Title) MaxLength() int {
	return constants.TITLE_MAX_LENGTH
}

func (a *Title) IsSameAs(b *Title) bool {
	return a.String() == b.String()
}

func (s *Title) MarshalBinary() ([]byte, error) {
	return MarshalStringToBytes(s.String(), s.MaxLength())
}

func (s *Title) UnmarshalBinary(data []byte) error {
	str, err := UnmarshalStringFromBytes(data, s.MaxLength())
	if err != nil {
		return err
	}
	err = s.SetString(str)
	if err != nil {
		return err
	}

	return nil
}

func (s *Title) UnmarshalBinaryData(data []byte) ([]byte, error) {
	str, data, err := UnmarshalStringFromBytesData(data, s.MaxLength())
	if err != nil {
		return data, err
	}

	err = s.SetString(str)
	if err != nil {
		return data, err
	}

	return data, nil
}
