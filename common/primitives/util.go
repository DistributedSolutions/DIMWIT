package primitives

import (
	"fmt"
)

func MarshalStringToBytes(str string, maxlength int) ([]byte, error) {
	if len(str) > maxlength {
		return nil, fmt.Errorf("Length of string is too long, found length is %d, max length is %d",
			len(str), maxlength)
	}

	data := []byte(str)
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
