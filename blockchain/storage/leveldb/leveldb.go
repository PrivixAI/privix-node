package leveldb

import (
	"fmt"

	"github.com/hashicorp/go-hclog"
	"github.com/PrivixAI-labs/Privix-node/blockchain/storage"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

const (
	DefaultCache   = int(256)
	DefaultHandles = int(256)
)

// Factory creates a leveldb storage
func Factory(config map[string]interface{}, logger hclog.Logger) (storage.Storage, error) {
	path, ok := config["path"]
	if !ok {
		return nil, fmt.Errorf("path not found")
	}

	pathStr, ok := path.(string)
	if !ok {
		return nil, fmt.Errorf("path is not a string")
	}

	return NewLevelDBStorage(pathStr, logger)
}

// NewLevelDBStorage creates the new storage reference with leveldb default options
func NewLevelDBStorage(path string, logger hclog.Logger) (storage.Storage, error) {
	// Set default options
	options := &opt.Options{
		OpenFilesCacheCapacity: DefaultHandles,
		BlockCacheCapacity:     DefaultCache / 2 * opt.MiB,
		WriteBuffer:            DefaultCache / 4 * opt.MiB, // Two of these are used internally
	}

	return NewLevelDBStorageWithOpt(path, logger, options)
}

// NewLevelDBStorageWithOpt creates the new storage reference with leveldb with custom options
func NewLevelDBStorageWithOpt(path string, logger hclog.Logger, opts *opt.Options) (storage.Storage, error) {
	db, err := leveldb.OpenFile(path, opts)
	if err != nil {
		return nil, err
	}

	kv := &levelDBKV{db}

	return storage.NewKeyValueStorage(logger.Named("leveldb"), kv), nil
}

// levelDBKV is the leveldb implementation of the kv storage
type levelDBKV struct {
	db *leveldb.DB
}

// Set sets the key-value pair in leveldb storage
func (l *levelDBKV) Set(p []byte, v []byte) error {
	return l.db.Put(p, v, &opt.WriteOptions{Sync: true})
}

// Get retrieves the key-value pair in leveldb storage
func (l *levelDBKV) Get(p []byte) ([]byte, bool, error) {
	data, err := l.db.Get(p, nil)
	if err != nil {
		if err.Error() == "leveldb: not found" {
			return nil, false, nil
		}

		return nil, false, err
	}

	return data, true, nil
}

// Close closes the leveldb storage instance
func (l *levelDBKV) Close() error {
	return l.db.Close()
}

func (l *levelDBKV) NewBatch() storage.Batch {
	return NewBatchLevelDB(l.db)
}
