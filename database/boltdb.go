package database

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/DistributedSolutions/DIMWIT/common/constants"
	"github.com/boltdb/bolt"
)

type BoltDB struct {
	db *bolt.DB

	sync.RWMutex
}

func NewBoltDB(filename string) IDatabase {
	db := new(BoltDB)

	if db.db == nil {
		if filename == "" {
			filename = "/tmp/bolt_my.db"
		}

		ss := strings.Split(filename, "/")

		if len(ss) > 1 {
			dir := ""
			ss = ss[:len(ss)-1]
			for _, s := range ss {
				dir += s + "/"
			}
			err := os.MkdirAll(dir, constants.DIRECTORY_PERMISSIONS)
			if err != nil {
				panic("Database was not found, and could not be created.")
			}
		}

		tdb, err := bolt.Open(filename, constants.FILE_PERMISSIONS, nil)
		if err != nil {
			panic("Database was not found, and could not be created.")
		}

		db.db = tdb
	}

	return db
}

func (db *BoltDB) ListAllBuckets() ([][]byte, error) {
	db.RLock()
	defer db.RUnlock()

	answer := [][]byte{}
	db.db.View(func(tx *bolt.Tx) error {
		c := tx.Cursor()
		k, _ := c.First()
		for {
			if k == nil {
				break
			}
			answer = append(answer, k)
			k, _ = c.Next()
		}
		return nil
	})

	return answer, nil
}

func (db *BoltDB) Delete(bucket []byte, key []byte) error {
	db.Lock()
	defer db.Unlock()

	db.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(bucket)
		if err != nil {
			return err
		}
		b := tx.Bucket(bucket)
		b.Delete(key)
		return nil
	})
	return nil
}

func (db *BoltDB) Close() error {
	db.Lock()
	defer db.Unlock()

	db.db.Close()
	return nil
}

func (db *BoltDB) Get(bucket []byte, key []byte) ([]byte, error) {
	db.RLock()
	defer db.RUnlock()

	var v []byte
	db.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		if b == nil {
			return nil
		}
		v = b.Get(key)
		if v == nil {
			return nil
		}
		return nil
	})
	if v == nil { // If the value is undefined, return nil
		return nil, nil
	}

	return v, nil
}

func (db *BoltDB) Put(bucket []byte, key []byte, data []byte) error {
	db.Lock()
	defer db.Unlock()

	err := db.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(bucket)
		if err != nil {
			return err
		}
		b := tx.Bucket(bucket)
		err = b.Put(key, data)
		return err
	})
	return err
}

type Record struct {
	Bucket []byte
	Key    []byte
	Data   []byte
}

func (db *BoltDB) PutInBatch(records []Record) error {
	db.Lock()
	defer db.Unlock()

	err := db.db.Batch(func(tx *bolt.Tx) error {
		for _, v := range records {
			_, err := tx.CreateBucketIfNotExists(v.Bucket)
			if err != nil {
				return err
			}
			b := tx.Bucket(v.Bucket)
			err = b.Put(v.Key, v.Data)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (db *BoltDB) Clear(bucket []byte) error {
	db.Lock()
	defer db.Unlock()

	err := db.db.Update(func(tx *bolt.Tx) error {
		err := tx.DeleteBucket(bucket)
		if err != nil {
			return fmt.Errorf("No bucket: %s", err)
		}
		return nil
	})
	return err
}

func (db *BoltDB) ListAllKeys(bucket []byte) (keys [][]byte, err error) {
	db.RLock()
	defer db.RUnlock()

	keys = make([][]byte, 0, 32)
	db.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			//fmt.Println("bucket 0x" + hex.EncodeToString(bucket) + " not found")
		} else {
			b.ForEach(func(k, v []byte) error {
				keys = append(keys, k)
				return nil
			})
		}
		return nil
	})
	return
}

func (db *BoltDB) DoesKeyExist(bucket, key []byte) (bool, error) {
	db.RLock()
	defer db.RUnlock()

	var v []byte
	db.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		if b == nil {
			return nil
		}
		v = b.Get(key)
		if v == nil {
			return nil
		}
		return nil
	})
	if v == nil { // If the value is undefined, return nil
		return false, nil
	}

	return true, nil
}

func (db *BoltDB) GetAll(bucket []byte) (data [][]byte, keys [][]byte, err error) {
	db.Lock()
	defer db.Unlock()

	data = [][]byte{}
	keys = [][]byte{}
	err = db.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			//fmt.Println("bucket 0x" + hex.EncodeToString(bucket) + " not found")
		} else {
			b.ForEach(func(k, v []byte) error {
				keys = append(keys, k)
				data = append(data, v)
				return nil
			})
			return nil
		}
		return nil
	})
	if err != nil {
		return nil, nil, err
	}
	return data, keys, nil
}
