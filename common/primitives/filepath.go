package primitives

import (
	"fmt"

	"github.com/DistributedSolutions/DIMWIT/common/constants"
	"github.com/DistributedSolutions/DIMWIT/common/primitives/random"
)

type FilePath string

func NewFilePath(url string) (*FilePath, error) {
	d := new(FilePath)

	err := d.SetString(url)
	if err != nil {
		return nil, err
	}

	return d, nil
}

func RandomFilePath() *FilePath {
	l, _ := NewFilePath("")
	l.SetString(random.RandStringOfSize(l.MaxLength()))

	return l
}

func (d *FilePath) Empty() bool {
	if d.String() == "" {
		return true
	}
	return false
}

func (d *FilePath) SetString(url string) error {
	if len(url) > d.MaxLength() {
		return fmt.Errorf("Description given is too long, length must be under %d, given length is %d",
			d.MaxLength(), len(url))
	}

	*d = FilePath(url)
	return nil
}

func (d *FilePath) String() string {
	return string(*d)
}

func (d *FilePath) MaxLength() int {
	return constants.FILE_PATH_MAX_LENGTH
}

func (a *FilePath) IsSameAs(b *FilePath) bool {
	return a.String() == b.String()
}

func (s *FilePath) MarshalBinary() ([]byte, error) {
	return MarshalStringToBytes(s.String(), s.MaxLength())
}

func (s *FilePath) UnmarshalBinary(data []byte) error {
	_, err := s.UnmarshalBinaryData(data)
	return err
}

func (s *FilePath) UnmarshalBinaryData(data []byte) (newData []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("[FilePath] A panic has occurred while unmarshaling: %s", r)
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
