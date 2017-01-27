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
	if len(title) > constants.TITLE_MAX_LENGTH {
		return fmt.Errorf("Description given is too long, length must be under %d, given length is %d",
			constants.TITLE_MAX_LENGTH, len(title))
	}

	*d = Title(title)
	return nil
}

func (d *Title) String() string {
	return string(*d)
}
