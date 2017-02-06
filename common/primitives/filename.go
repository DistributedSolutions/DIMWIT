package primitives

import (
	"bytes"
	"fmt"

	"github.com/DistributedSolutions/DIMWIT/common/constants"
	"github.com/DistributedSolutions/DIMWIT/common/primitives/random"
	"github.com/dustin/go-humanize"
)

type FileList struct {
	length   uint32
	fileList []File
}

func NewFileList() *FileList {
	af := new(FileList)
	af.fileList = make([]File, 0)
	af.length = 0

	return af
}

func RandomFileList(max uint32) *FileList {
	fl := NewFileList()
	l := random.RandomUInt32Between(0, max)
	fl.fileList = make([]File, l)

	for i := range fl.fileList {
		fl.fileList[i] = *(RandomFile())
	}

	fl.length = l
	return fl
}

func (af *FileList) AddFile(filename string, size int64) error {
	f, err := NewFile(filename, size)
	if err != nil {
		return err
	}

	af.length++
	af.fileList = append(af.fileList, *f)
	return nil
}

func (fl *FileList) GetFiles() []File {
	return fl.fileList
}

func (a *FileList) IsSameAs(b *FileList) bool {
	if a.length != b.length {
		return false
	}

	for i := range a.fileList {
		if !a.fileList[i].IsSameAs(&(b.fileList[i])) {
			return false
		}
	}

	return true
}

func (fl *FileList) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	data := Uint32ToBytes(fl.length)
	buf.Write(data)

	for _, f := range fl.fileList {
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
	fl.length = u

	fl.fileList = make([]File, fl.length)

	var i uint32 = 0
	for ; i < fl.length; i++ {
		newData, err = fl.fileList[i].UnmarshalBinaryData(newData)
		if err != nil {
			return data, err
		}
	}

	return
}

type File struct {
	fileName string // includes extension
	size     int64
	checksum MD5Checksum
}

func NewFile(filename string, size int64) (*File, error) {
	f := new(File)

	err := f.SetFileName(filename)
	if err != nil {
		return nil, err
	}

	f.size = size

	return f, nil
}

func RandomFile() *File {
	f := new(File)
	str := random.RandStringOfSize(f.MaxLength())
	f.SetFileName(str)
	s := random.RandomInt64()
	f.SetSize(s)
	return f

}

/*func (f *File) GetFileName() string {
	return f.FileName
}*/

func (f *File) GetFullPath() string {
	return f.fileName
}

func (f *File) SetFileName(filename string) error {
	if len(filename) > constants.FILE_NAME_MAX_LENGTH {
		return fmt.Errorf("Name given is too long, length must be under %d, given length is %d", constants.FILE_NAME_MAX_LENGTH, len(filename))
	}

	f.fileName = filename
	return nil
}

func (f *File) SetSize(size int64) {
	f.size = size
}

func (f *File) GetSize() int64 {
	return f.size
}

func (f *File) String() string {
	return fmt.Sprintf("%s (%s)", f.fileName, humanize.Bytes(uint64(f.size)))
}

func (d *File) MaxLength() int {
	return constants.FILE_NAME_MAX_LENGTH
}

func (a *File) IsSameAs(b *File) bool {
	if a.fileName != b.fileName {
		return false
	}

	if a.size != b.size {
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

	data, err = Int64ToBytes(f.size)
	if err != nil {
		return nil, err
	}
	buf.Write(data)

	data, err = f.checksum.MarshalBinary()
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

	newData, err = f.checksum.UnmarshalBinaryData(newData)
	if err != nil {
		return data, err
	}

	return
}
