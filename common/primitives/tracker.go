package primitives

import (
	"bytes"
	"fmt"

	"github.com/DistributedSolutions/DIMWIT/common/constants"
	"github.com/DistributedSolutions/DIMWIT/common/primitives/random"
)

type TrackerList struct {
	Trackers []Tracker `json:"trackers"`
}

func NewTrackerList() *TrackerList {
	tl := new(TrackerList)
	tl.Trackers = make([]Tracker, 0)

	return tl
}

func RandomTrackerList(max uint32) *TrackerList {
	tl := NewTrackerList()
	l := random.RandomUInt32Between(0, max)

	var i uint32
	for i = 0; i < l; i++ {
		tl.AddNewTracker(RandomTracker())
	}

	return tl
}

func (d *TrackerList) Empty() bool {
	if len(d.Trackers) == 0 {
		return true
	}
	return false
}

func (tl *TrackerList) AddNewTracker(tracker *Tracker) error {
	tl.Trackers = append(tl.Trackers, *tracker)

	return nil
}

func (tl *TrackerList) AddNewTrackerURL(url string) error {
	t, err := NewTracker(url)
	if err != nil {
		return err
	}

	tl.Trackers = append(tl.Trackers, *t)

	return nil
}

func (tl *TrackerList) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	data := Uint32ToBytes(uint32(len(tl.Trackers)))
	buf.Write(data)

	for _, t := range tl.Trackers {
		data, err := t.MarshalBinary()
		if err != nil {
			return nil, err
		}
		buf.Write(data)
	}

	return buf.Next(buf.Len()), nil
}

func (tl *TrackerList) UnmarshalBinary(data []byte) error {
	_, err := tl.UnmarshalBinaryData(data)
	return err
}

func (tl *TrackerList) UnmarshalBinaryData(data []byte) (newData []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("[TrackerList] A panic has occurred while unmarshaling: %s", r)
			return
		}
	}()

	newData = data

	u, err := BytesToUint32(newData[:4])
	if err != nil {
		return data, err
	}
	newData = newData[4:]

	tl.Trackers = make([]Tracker, u)
	var i uint32
	for i = 0; i < u; i++ {
		newData, err = tl.Trackers[i].UnmarshalBinaryData(newData)
		if err != nil {
			return data, err
		}
	}

	return
}

func (a *TrackerList) IsSameAs(b *TrackerList) bool {
	if len(a.Trackers) != len(b.Trackers) {
		return false
	}

	for i := range a.Trackers {
		if !a.Trackers[i].IsSameAs(&(b.Trackers[i])) {
			return false
		}
	}

	return true
}

type Tracker string

func NewTracker(url string) (*Tracker, error) {
	d := new(Tracker)

	err := d.SetString(url)
	if err != nil {
		return nil, err
	}

	return d, nil
}

func (d *Tracker) Empty() bool {
	if d.String() == "" {
		return true
	}
	return false
}

func (t *Tracker) SetString(url string) error {
	if len(url) > t.MaxLength() {
		return fmt.Errorf("Tracker url given is too long, length must be under %d, given length is %d",
			t.MaxLength(), len(url))
	}

	*t = Tracker(url)
	return nil
}

func (t *Tracker) String() string {
	return string(*t)
}

func (d *Tracker) MaxLength() int {
	return constants.TRACKER_URL_MAX_LENGTH
}

func (a *Tracker) IsSameAs(b *Tracker) bool {
	return a.String() == b.String()
}

func (d *Tracker) MarshalBinary() ([]byte, error) {
	return MarshalStringToBytes(d.String(), d.MaxLength())
}

func (t *Tracker) UnmarshalBinary(data []byte) error {
	_, err := t.UnmarshalBinaryData(data)
	return err
}

func (d *Tracker) UnmarshalBinaryData(data []byte) (newData []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("[Tracker] A panic has occurred while unmarshaling: %s", r)
			return
		}
	}()

	newData = data
	str, newData, err := UnmarshalStringFromBytesData(newData, d.MaxLength())
	if err != nil {
		return data, err
	}

	err = d.SetString(str)
	if err != nil {
		return data, err
	}

	return newData, nil
}

func RandomTracker() *Tracker {
	l, _ := NewTracker("")
	l.SetString(random.RandStringOfSize(l.MaxLength()))

	return l
}
