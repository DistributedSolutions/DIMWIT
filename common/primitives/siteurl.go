package primitives

import (
	"fmt"

	"github.com/DistributedSolutions/DIMWIT/common/constants"
)

type SiteURL string

func NewURL(url string) (*SiteURL, error) {
	d := new(SiteURL)

	err := d.SetString(url)
	if err != nil {
		return nil, err
	}

	return d, nil
}

func (d *SiteURL) SetString(url string) error {
	if len(url) > constants.URL_MAX_LENGTH {
		return fmt.Errorf("Description given is too long, length must be under %d, given length is %d",
			constants.URL_MAX_LENGTH, len(url))
	}

	*d = SiteURL(url)
	return nil
}

func (d *SiteURL) String() string {
	return string(*d)
}
