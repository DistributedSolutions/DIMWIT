package primitives

import (
	"bytes"
	"fmt"

	"github.com/DistributedSolutions/DIMWIT/common/constants"
	"github.com/DistributedSolutions/DIMWIT/common/primitives/random"
	"github.com/dustin/go-humanize"
)

type FileList struct {
	FileList []File `json:"filelist"`
}

func NewFileList() *FileList {
	af := new(FileList)
	af.FileList = make([]File, 0)

	return af
}

func RandomFileList(max uint32) *FileList {
	fl := NewFileList()
	l := random.RandomUInt32Between(0, max)
	fl.FileList = make([]File, l)

	for i := range fl.FileList {
		fl.FileList[i] = *(RandomFile())
	}
	return fl
}

func (af *FileList) Empty() bool {
	if len(af.FileList) == 0 {
		return true
	}

	return false
}

func (af *FileList) AddFile(filename string, size int64) error {
	f, err := NewFile(filename, size)
	if err != nil {
		return err
	}

	af.FileList = append(af.FileList, *f)
	return nil
}

func (fl *FileList) GetFiles() []File {
	return fl.FileList
}

func (a *FileList) IsSameAs(b *FileList) bool {
	if len(a.FileList) != len(b.FileList) {
		return false
	}

	for i := range a.FileList {
		if !a.FileList[i].IsSameAs(&(b.FileList[i])) {
			return false
		}
	}

	return true
}

func (fl *FileList) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	data := Uint32ToBytes(uint32(len(fl.FileList)))
	buf.Write(data)

	for _, f := range fl.FileList {
		data, err := f.MarshalBinary()
		if err != nil {
			return nil, err
		}
		buf.Write(data)
	}

	return buf.Next(buf.Len()), nil
}

func (fl *FileList) UnmarshalBinary(data []byte) error {
	_, err := fl.UnmarshalBinaryData(data)
	return err
}

func (fl *FileList) UnmarshalBinaryData(data []byte) (newData []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("[FileName] A panic has occurred while unmarshaling: %s", r)
			return
		}
	}()

	newData = data

	u, err := BytesToUint32(newData[:4])
	if err != nil {
		return data, err
	}

	newData = newData[4:]

	fl.FileList = make([]File, u)

	var i uint32 = 0
	for ; i < u; i++ {
		newData, err = fl.FileList[i].UnmarshalBinaryData(newData)
		if err != nil {
			return data, err
		}
	}

	return
}

type File struct {
	FileName string      `json:"file"` // includes extension
	Size     int64       `json:"size"`
	Checksum MD5Checksum `json:"checksum"`
}

func NewFile(filename string, size int64) (*File, error) {
	f := new(File)

	err := f.SetFileName(filename)
	if err != nil {
		return nil, err
	}

	f.Size = size

	return f, nil
}

func RandomFile() *File {
	f := new(File)
	str := random.RandStringOfSize(f.MaxLength())
	f.SetFileName(str)
	s := random.RandomInt64()
	f.SetSize(s)
	f.Checksum = *RandomMD5()
	return f

}

func (af *File) Empty() bool {
	if af.FileName == "" || af.Checksum.Empty() {
		return true
	}

	return false
}

/*func (f *File) GetFileName() string {
	return f.FileName
}*/

func (f *File) GetFullPath() string {
	return f.FileName
}

func (f *File) SetFileName(filename string) error {
	if len(filename) > constants.FILE_NAME_MAX_LENGTH {
		return fmt.Errorf("Name given is too long, length must be under %d, given length is %d", constants.FILE_NAME_MAX_LENGTH, len(filename))
	}

	f.FileName = filename
	return nil
}

func (f *File) SetSize(size int64) {
	f.Size = size
}

func (f *File) GetSize() int64 {
	return f.Size
}

func (f *File) String() string {
	return fmt.Sprintf("%s (%s)", f.FileName, humanize.Bytes(uint64(f.Size)))
}

func (d *File) MaxLength() int {
	return constants.FILE_NAME_MAX_LENGTH
}

func (a *File) IsSameAs(b *File) bool {
	if a.FileName != b.FileName {
		return false
	}

	if a.Size != b.Size {
		return false
	}

	return true
}

func (f *File) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	data, err := MarshalStringToBytes(f.GetFullPath(), f.MaxLength())
	if err != nil {
		return nil, err
	}
	buf.Write(data)

	data, err = Int64ToBytes(f.Size)
	if err != nil {
		return nil, err
	}
	buf.Write(data)

	data, err = f.Checksum.MarshalBinary()
	if err != nil {
		return nil, err
	}
	buf.Write(data)

	return buf.Next(buf.Len()), nil
}

func (f *File) UnmarshalBinary(data []byte) error {
	_, err := f.UnmarshalBinaryData(data)
	return err
}

func (f *File) UnmarshalBinaryData(data []byte) (newData []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("A panic has occurred while unmarshaling: %s", r)
			return
		}
	}()

	newData = data

	str, newData, err := UnmarshalStringFromBytesData(newData, f.MaxLength())
	if err != nil {
		return data, err
	}

	err = f.SetFileName(str)
	if err != nil {
		return data, err
	}

	val, err := BytesToInt64(newData[:8])
	if err != nil {
		return data, err
	}
	f.SetSize(val)

	newData = newData[8:]

	newData, err = f.Checksum.UnmarshalBinaryData(newData)
	if err != nil {
		return data, err
	}

	return
}
