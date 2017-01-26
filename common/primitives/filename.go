package primitives

import (
	"fmt"

	"github.com/DistributedSolutions/DIMWIT/common/constants"
)

type FileList struct {
	Length   int
	FileList []File
}

func NewFileList() *FileList {
	af := new(FileList)
	af.FileList = make([]File, 0)
	af.Length = 0

	return af
}

func (af *FileList) AddFile(filename string, extension string, size int64) error {
	f, err := NewFile(filename, extension, size)
	if err != nil {
		return err
	}

	af.Length++
	af.FileList = append(af.FileList, *f)
	return nil
}

type File struct {
	FileName  string
	Extension string
	Size      int64
}

func NewFile(filename string, extension string, size int64) (*File, error) {
	f := new(File)

	err := f.SetFileName(filename)
	if err != nil {
		return nil, err
	}

	err = f.SetExtension(extension)
	if err != nil {
		return nil, err
	}

	f.Size = size

	return f, nil
}

func (f *File) GetFileName() string {
	return f.FileName
}

func (f *File) GetExtensions() string {
	return f.Extension
}

func (f *File) GetFullPath() string {
	return f.FileName + f.Extension
}

func (f *File) SetFileName(filename string) error {
	if len(filename) > constants.FILE_NAME_MAX_LENGTH {
		return fmt.Errorf("Name given is too long, length must be under %d, given length is %d", constants.FILE_NAME_MAX_LENGTH, len(filename))
	}

	f.FileName = filename
	return nil
}

func (f *File) SetExtension(extension string) error {
	if len(extension) > constants.FILENAME_EXTENSION_MAX_LENGTH {
		return fmt.Errorf("Extension given is too long, length must be under %d, given length is %d", constants.FILENAME_EXTENSION_MAX_LENGTH, len(extension))
	}

	f.Extension = extension
	return nil
}

func (f *File) SetSize(size int64) {
	f.Size = size
}

func (f *File) GetSize() int64 {
	return f.Size
}
