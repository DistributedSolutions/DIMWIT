package database

type IDatabase interface {
	ListAllBuckets() ([][]byte, error)
	Delete(bucket []byte, key []byte) error
	Close() error
	Get(bucket []byte, key []byte) ([]byte, error)
	Put(bucket []byte, key []byte, data []byte) error
	PutInBatch(records []Record) error
	Clear(bucket []byte) error
	ListAllKeys(bucket []byte) (keys [][]byte, err error)
	DoesKeyExist(bucket, key []byte) (bool, error)
	GetAll(bucket []byte) (data [][]byte, keys [][]byte, err error)
}
