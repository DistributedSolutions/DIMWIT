package primitives

import (
	"fmt"

	"github.com/DistributedSolutions/DIMWIT/common/constants"
	"github.com/DistributedSolutions/DIMWIT/common/primitives/random"
)

type Title string

func RandomTitle() *Title {
	l, _ := NewTitle("")
	l.SetString(random.RandStringOfSize(l.MaxLength()))

	return l
}

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
		return fmt.Errorf("Title given is too long, length must be under %d, given length is %d",
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

func (t *Title) UnmarshalBinary(data []byte) error {
	_, err := t.UnmarshalBinaryData(data)
	return err
}

func (t *Title) UnmarshalBinaryData(data []byte) (newData []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("[Title] A panic has occurred while unmarshaling: %s", r)
			return
		}
	}()

	newData = data
	str, newData, err := UnmarshalStringFromBytesData(newData, t.MaxLength())
	if err != nil {
		return data, err
	}

	err = t.SetString(str)
	if err != nil {
		return data, err
	}

	return newData, nil
}
