package database_test

import (
	"bytes"
	"testing"

	"github.com/DistributedSolutions/DIMWIT/common/primitives/random"
	. "github.com/DistributedSolutions/DIMWIT/database"
)

func TestMap(t *testing.T) {
	return
	database := NewMapDB()
	defer database.Close()

	var amt int = 1000
	keys := make([][]byte, amt)
	datas := make([][]byte, amt)
	for i := 0; i < amt; i++ {
		keys[i] = random.RandByteSliceOfSize(32)
		datas[i] = random.RandByteSliceOfSize(random.RandomIntBetween(0, 10000))
	}

	for i := 0; i < amt/2; i++ {
		err := database.Put([]byte("singles"), keys[i], datas[i])
		if err != nil {
			t.Error(err)
		}
	}

	for i := amt / 2; i < amt; i += 2 {
		r := make([]Record, 2)
		r[0].Bucket = []byte("doubles")
		r[0].Data = datas[i]
		r[0].Key = keys[i]

		r[1].Bucket = []byte("doubles")
		r[1].Data = datas[i+1]
		r[1].Key = keys[i+1]
		err := database.PutInBatch(r)
		if err != nil {
			t.Error(err)
		}
	}

	for i := 0; i < amt/2; i++ {
		data, err := database.Get([]byte("singles"), keys[i])
		if err != nil {
			t.Error(err)
		} else if bytes.Compare(data, datas[i]) != 0 {
			t.Error("Wrong data back")
		}
	}

	for i := amt / 2; i < amt; i++ {
		data, err := database.Get([]byte("doubles"), keys[i])
		if err != nil {
			t.Error(err)
		} else if bytes.Compare(data, datas[i]) != 0 {
			t.Error("Wrong data back")
		}
	}

	for i := 0; i < amt/10; i++ {
		err := database.Put([]byte("clearthis"), keys[i], datas[i])
		if err != nil {
			t.Error(err)
		}
	}

	buckeys, err := database.ListAllKeys([]byte("clearthis"))
	if err != nil {
		t.Error(err)
	}

	hold := make(map[string]bool)
	for _, k := range buckeys {
		hold[string(k)] = true
	}

	for i := 0; i < amt/10; i++ {
		if _, ok := hold[string(keys[i])]; !ok {
			t.Error("Key was not listed")
		}
	}

	bucs, err := database.ListAllBuckets()
	if err != nil {
		t.Error(err)
	}

	if len(bucs) != 3 {
		t.Error("Not all buckets fetched")
	}

	err = database.Clear([]byte("clearthis"))
	if err != nil {
		t.Error(err)
	}

	buckeys, err = database.ListAllKeys([]byte("clearthis"))
	if err != nil {
		t.Error(err)
	}
	if len(buckeys) != 0 {
		t.Error("Should be 0")
	}

	for i := 0; i < amt/2; i++ {
		e, err := database.DoesKeyExist([]byte("singles"), keys[i])
		if err != nil {
			t.Error(err)
		}
		if !e {
			t.Error("Should exist")
		}
	}

	for i := 0; i < amt/2; i++ {
		err := database.Delete([]byte("singles"), keys[i])
		if err != nil {
			t.Error(err)
		}
	}

	for i := 0; i < amt/2; i++ {
		e, err := database.DoesKeyExist([]byte("singles"), keys[i])
		if err != nil {
			t.Error(err)
		}
		if e {
			t.Error("Should not exist")
		}
	}

	bucDatas, bucKeys, err := database.GetAll([]byte("doubles"))
	if err != nil {
		t.Error(err)
	}

	if len(bucDatas) != amt/2 {
		t.Errorf("Bad length, length should be %d, it is %d", amt/2, len(bucDatas))
	}
	if len(bucKeys) != amt/2 {
		t.Errorf("Bad length, length should be %d, it is %d", amt/2, len(bucKeys))
	}
}
