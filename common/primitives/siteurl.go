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
	if len(url) > d.MaxLength() {
		return fmt.Errorf("Description given is too long, length must be under %d, given length is %d",
			d.MaxLength(), len(url))
	}

	*d = SiteURL(url)
	return nil
}

func (d *SiteURL) String() string {
	return string(*d)
}

func (d *SiteURL) MaxLength() int {
	return constants.URL_MAX_LENGTH
}

func (a *SiteURL) IsSameAs(b *SiteURL) bool {
	return a.String() == b.String()
}

func (d *SiteURL) MarshalBinary() ([]byte, error) {
	return MarshalStringToBytes(d.String(), d.MaxLength())
}

func (d *SiteURL) UnmarshalBinary(data []byte) error {
	str, err := UnmarshalStringFromBytes(data, d.MaxLength())
	if err != nil {
		return err
	}
	err = d.SetString(str)
	if err != nil {
		return err
	}

	return nil
}

func (d *SiteURL) UnmarshalBinaryData(data []byte) ([]byte, error) {
	str, data, err := UnmarshalStringFromBytesData(data, d.MaxLength())
	if err != nil {
		return data, err
	}

	err = d.SetString(str)
	if err != nil {
		return data, err
	}

	return data, nil
}
