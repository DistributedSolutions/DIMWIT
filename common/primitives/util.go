package primitives

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"fmt"
)

func Float64ToBytes(v float64) ([]byte, error) {
	str := strconv.FormatFloat(v, 'f', 6, 64)
	b, err := MarshalStringToBytes(str, 200)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func BytesToFloat64(data []byte) (float64, error) {
	v, _, e := BytesToFloat64Data(data)
	return v, e
}

func BytesToFloat64Data(data []byte) (float64, []byte, error) {
	var err error
	var resp string
	newData := data

	resp, newData, err = UnmarshalStringFromBytesData(newData, 200)
	if err != nil {
		return 0, data, err
	}
	f, err := strconv.ParseFloat(resp, 64)
	if err != nil {
		return 0, data, err
	}

	return f, newData, nil
}

func UnmarshalBinarySlice(data []byte) (resp []byte, err error) {
	r, _, e := UnmarshalBinarySliceData(data)
	return r, e
}

func UnmarshalBinarySliceData(data []byte) (resp []byte, newData []byte, err error) {
	newData = data

	u, err := BytesToUint32(data)
	if err != nil {
		return nil, data, err
	}

	newData = newData[4:]
	if len(newData) < int(u) {
		return nil, data, fmt.Errorf("need at least %d bytes, found %d", u+4, len(data))
	}

	resp = newData[:u]
	newData = newData[u:]
	return
}

func MarshalBinarySlice(slice []byte) []byte {
	buf := new(bytes.Buffer)

	u := uint32(len(slice))
	buf.Write(Uint32ToBytes(u))
	buf.Write(slice)

	return buf.Next(buf.Len())
}

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
			break
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

func Uint32ToBytes(val uint32) []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, val)

	return b
}

func BytesToUint32(data []byte) (ret uint32, err error) {
	buf := bytes.NewBuffer(data)
	err = binary.Read(buf, binary.BigEndian, &ret)
	return
}

func Uint64ToBytes(val uint64) ([]byte, error) {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, val)

	return b, nil
}

func BytesToUint64(data []byte) (uint64, error) {
	if len(data) != 8 {
		return 0, fmt.Errorf("Must be exactly 8 bytes in length, found length of %d bytes", len(data))
	}
	u := binary.BigEndian.Uint64(data)
	return u, nil
}

func Int64ToBytes(val int64) ([]byte, error) {
	return Uint64ToBytes(uint64(val))
}

func BytesToInt64(data []byte) (int64, error) {
	val, err := BytesToUint64(data)
	return int64(val), err
}

func BoolToBytes(b bool) []byte {
	if b {
		return []byte{0x01}
	}
	return []byte{0x00}
}

func ByteToBool(b byte) bool {
	if b == 0x00 {
		return false
	}
	return true
}

func RandXORKey() byte {
	xorCipher := make([]byte, 1)
	rand.Read(xorCipher)
	if xorCipher[0] == 0x00 {
		return RandXORKey()
	} else {
		return xorCipher[0]
	}
}

func XORCipher(key byte, data []byte) []byte {
	buf := new(bytes.Buffer)

	for _, d := range data {
		buf.Write([]byte{d ^ key})
	}

	return buf.Next(buf.Len())
}
