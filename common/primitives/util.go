package primitives

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

func MarshalStringToBytes(str string, maxlength int) ([]byte, error) {
	if len(str) > maxlength {
		return nil, fmt.Errorf("Length of string is too long, found length is %d, max length is %d",
			len(str), maxlength)
	}

	data := []byte(str)
	for i := 0; i < len(data); i++ {
		if data[i] == 0x00 {
			// Naughty, Naughty, Naughty
			data[i] = 0x01
		}
	}
	data = append(data, 0x00)

	return data, nil
}

func UnmarshalStringFromBytes(data []byte, maxlength int) (resp string, err error) {
	resp, _, err = UnmarshalStringFromBytesData(data, maxlength)
	return
}

func UnmarshalStringFromBytesData(data []byte, maxlength int) (resp string, newData []byte, err error) {
	newData = data
	end := -1
	if len(data)-1 < maxlength {
		maxlength = len(data) - 1
	}

	for i := 0; i <= maxlength; i++ {
		if newData[i] == 0x00 {
			// found null terminator
			end = i
		}
	}

	if end == -1 {
		err = fmt.Errorf("Could not find a 0x00 byte before max length + 1")
		return
	}

	resp = string(newData[:end])
	newData = newData[end+1:]
	return
}

func Uint32ToBytes(val uint32) ([]byte, error) {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, val)

	return b, nil
}

func BytesToUInt32(data []byte) (ret uint32, err error) {
	buf := bytes.NewBuffer(data)
	err = binary.Read(buf, binary.LittleEndian, &ret)
	return
}
