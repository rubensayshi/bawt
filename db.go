package bawt

import (
	"encoding/json"
	"fmt"

	"github.com/boltdb/bolt"
)

const bawtDBDefaultBucket = "bawt"

// GetDBKey retrieves a `key` from persistent storage and JSON
// unmarshales it into `v`.  We need to `Update` otherwise
// CreateBucketIfNotExists cannot create a bucket and returns
// an error immediately.
func (bot *Bot) GetDBKey(key string, v interface{}) error {
	return bot.DB.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(bawtDBDefaultBucket))
		if err != nil {
			return err
		}

		val := bucket.Get([]byte(key))
		if val == nil {
			return fmt.Errorf("not found")
		}

		err = json.Unmarshal(val, &v)
		if err != nil {
			return err
		}

		return nil
	})
}

// PutDBKey sets a key to the specified value in the persistent storage
// it JSON marshals the value before storing it.
func (bot *Bot) PutDBKey(key string, v interface{}) error {
	return bot.DB.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(bawtDBDefaultBucket))
		if err != nil {
			return err
		}

		jsonRes, err := json.Marshal(v)
		if err != nil {
			return err
		}

		err = bucket.Put([]byte(key), jsonRes)
		if err != nil {
			return err
		}

		return nil
	})
}
