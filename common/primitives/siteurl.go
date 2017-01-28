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

func (s *SiteURL) MarshalBinary() ([]byte, error) {
	return MarshalStringToBytes(s.String(), s.MaxLength())
}

func (s *SiteURL) UnmarshalBinary(data []byte) error {
	_, err := s.UnmarshalBinaryData(data)
	return err
}

func (s *SiteURL) UnmarshalBinaryData(data []byte) (newData []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("A panic has occurred while unmarshaling: %s", r)
			return
		}
	}()

	newData = data
	str, newData, err := UnmarshalStringFromBytesData(newData, s.MaxLength())
	if err != nil {
		return data, err
	}

	err = s.SetString(str)
	if err != nil {
		return data, err
	}

	return
}
