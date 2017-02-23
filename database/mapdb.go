package database

import (
	"sync"
)

type MapDB struct {
	Sem   sync.RWMutex
	Cache map[string]map[string][]byte // Our Cache
}

func NewMapDB() IDatabase {
	db := new(MapDB)
	db.Cache = make(map[string]map[string][]byte)
	return db
}

func (MapDB) Close() error {
	return nil
}

func (db *MapDB) ListAllBuckets() ([][]byte, error) {
	if db.Cache == nil {
		db.Sem.Lock()
		db.Cache = map[string]map[string][]byte{}
		db.Sem.Unlock()
	}

	db.Sem.RLock()
	defer db.Sem.RUnlock()

	answer := [][]byte{}
	for k, _ := range db.Cache {
		answer = append(answer, []byte(k))
	}

	return answer, nil
}

// Don't do anything here.
func (db *MapDB) Trim() {
}

func (db *MapDB) createCache(bucket []byte) {
	if db.Cache == nil {
		db.Sem.Lock()
		db.Cache = map[string]map[string][]byte{}
		db.Sem.Unlock()
	}
	db.Sem.RLock()
	_, ok := db.Cache[string(bucket)]
	db.Sem.RUnlock()
	if ok == false {
		db.Sem.Lock()
		db.Cache[string(bucket)] = map[string][]byte{}
		db.Sem.Unlock()
	}
}

func (db *MapDB) Init(bucketList [][]byte) {
	db.Sem.Lock()
	defer db.Sem.Unlock()

	db.Cache = map[string]map[string][]byte{}
	for _, v := range bucketList {
		db.Cache[string(v)] = map[string][]byte{}
	}
}

func (db *MapDB) Put(bucket, key []byte, data []byte) error {
	db.Sem.Lock()
	defer db.Sem.Unlock()

	if db.Cache == nil {
		db.Cache = map[string]map[string][]byte{}
	}
	_, ok := db.Cache[string(bucket)]
	if ok == false {
		db.Cache[string(bucket)] = map[string][]byte{}
	}

	db.Cache[string(bucket)][string(key)] = data
	return nil
}

func (db *MapDB) PutInBatch(records []Record) error {
	for _, v := range records {
		err := db.Put(v.Bucket, v.Key, v.Data)
		if err != nil {
			return err
		}
	}
	return nil
}

func (db *MapDB) Get(bucket, key []byte) ([]byte, error) {
	db.createCache(bucket)

	db.Sem.RLock()
	defer db.Sem.RUnlock()

	if db.Cache == nil {
		db.Cache = map[string]map[string][]byte{}
	}
	_, ok := db.Cache[string(bucket)]
	if ok == false {
		db.Cache[string(bucket)] = map[string][]byte{}
	}
	v, ok := db.Cache[string(bucket)][string(key)]
	if ok == false {
		return nil, nil
	}
	if v == nil {
		return nil, nil
	}
	return v, nil
}

func (db *MapDB) Delete(bucket, key []byte) error {
	db.Sem.Lock()
	defer db.Sem.Unlock()

	if db.Cache == nil {
		db.Cache = map[string]map[string][]byte{}
	}
	_, ok := db.Cache[string(bucket)]
	if ok == false {
		db.Cache[string(bucket)] = map[string][]byte{}
	}
	delete(db.Cache[string(bucket)], string(key))
	return nil
}

func (db *MapDB) ListAllKeys(bucket []byte) ([][]byte, error) {
	db.createCache(bucket)

	db.Sem.RLock()
	defer db.Sem.RUnlock()

	if db.Cache == nil {
		db.Cache = map[string]map[string][]byte{}
	}
	_, ok := db.Cache[string(bucket)]
	if ok == false {
		db.Cache[string(bucket)] = map[string][]byte{}
	}
	answer := [][]byte{}
	for k, _ := range db.Cache[string(bucket)] {
		answer = append(answer, []byte(k))
	}

	return answer, nil
}

func (db *MapDB) GetAll(bucket []byte) ([][]byte, [][]byte, error) {
	db.createCache(bucket)

	db.Sem.RLock()
	defer db.Sem.RUnlock()

	if db.Cache == nil {
		db.Cache = map[string]map[string][]byte{}
	}
	_, ok := db.Cache[string(bucket)]
	if ok == false {
		db.Cache[string(bucket)] = map[string][]byte{}
	}

	keys, err := db.ListAllKeys(bucket)
	if err != nil {
		return nil, nil, err
	}

	answer := make([][]byte, 0)
	for _, k := range keys {
		v := db.Cache[string(bucket)][string(k)]
		answer = append(answer, v)
	}
	return keys, answer, nil
}

func (db *MapDB) Clear(bucket []byte) error {
	db.Sem.Lock()
	defer db.Sem.Unlock()

	if db.Cache == nil {
		db.Cache = map[string]map[string][]byte{}
	}
	delete(db.Cache, string(bucket))
	return nil
}

func (db *MapDB) DoesKeyExist(bucket, key []byte) (bool, error) {
	db.createCache(bucket)

	db.Sem.RLock()
	defer db.Sem.RUnlock()

	if db.Cache == nil {
		db.Cache = map[string]map[string][]byte{}
	}
	_, ok := db.Cache[string(bucket)]
	if ok == false {
		db.Cache[string(bucket)] = map[string][]byte{}
	}
	_, ok = db.Cache[string(bucket)][string(key)]
	return ok, nil
}
