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
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("A panic has occurred while unmarshaling: %s", r)
		}
	}()

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

func BytesToUint32(data []byte) (ret uint32, err error) {
	buf := bytes.NewBuffer(data)
	err = binary.Read(buf, binary.LittleEndian, &ret)
	return
}

func Uint64ToBytes(val uint64) ([]byte, error) {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, val)

	return b, nil
}

func BytesToUint64(data []byte) (uint64, error) {
	if len(data) != 8 {
		return 0, fmt.Errorf("Must be exactly 8 bytes in length, found length of %d bytes", len(data))
	}
	u := binary.LittleEndian.Uint64(data)
	return u, nil
}

func Int64ToBytes(val int64) ([]byte, error) {
	return Uint64ToBytes(uint64(val))
}

func BytesToInt64(data []byte) (int64, error) {
	val, err := BytesToUint64(data)
	return int64(val), err
}

func BytesIsSame(a []byte, b []byte) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
