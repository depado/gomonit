package database

import (
	"fmt"
	"time"

	"github.com/boltdb/bolt"
)

// Storage is a type that contains a bolt.DB and a boolean that indicates if the connection is already open or not.
type Storage struct {
	DB     *bolt.DB
	Path   string
	Opened bool
}

// Open opens the database connection and create the file if necessary
func (s *Storage) Open() error {
	var err error
	dbfile := s.Path
	config := &bolt.Options{Timeout: 1 * time.Second}
	s.DB, err = bolt.Open(dbfile, 0600, config)
	if err == nil {
		s.Opened = true
	}
	return err
}

// Close closes the connection (or at least attempts to)
func (s *Storage) Close() error {
	s.Opened = false
	err := s.DB.Close()
	return err
}

// Storable is the type of data that can be stored/retrieved from the database.
type Storable interface {
	Encode() ([]byte, error)
	Decode([]byte) error
}

// Save saves some data inside the bucket at the specified key.
func (s Storage) Save(bucket, key string, data Storable) error {
	if !s.Opened {
		return fmt.Errorf("db must be opened before saving")
	}
	err := s.DB.Update(func(tx *bolt.Tx) error {
		mBucket, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return fmt.Errorf("Error creating bucket : %s", err)
		}
		enc, err := data.Encode()
		if err != nil {
			return fmt.Errorf("Could not encode : %s", err)
		}
		err = mBucket.Put([]byte(key), enc)
		return err
	})
	return err
}

// Delete deletes data inside the bucket at the specified key.
func (s Storage) Delete(bucket, key string) error {
	if !s.Opened {
		return fmt.Errorf("db must be opened before using it")
	}
	err := s.DB.Update(func(tx *bolt.Tx) error {
		mBucket := tx.Bucket([]byte(bucket))

		if mBucket != nil {
			err := mBucket.Delete([]byte(key))
			return err
		}
		return nil
	})
	return err
}

// Get retrieves the specific Storable object from bucket and key
func (s Storage) Get(bucket, key string, to Storable) error {
	if !s.Opened {
		return fmt.Errorf("Database must be opened first.")
	}
	err := s.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		k := []byte(key)
		err := to.Decode(b.Get(k))
		return err
	})
	return err
}

// List keys
func (s Storage) List(bucket string, to *[]string) error {
	if !s.Opened {
		return fmt.Errorf("Database must be opened first.")
	}
	err := s.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b != nil {
			err := b.ForEach(func(k, _ []byte) error {
				*to = append(*to, fmt.Sprintf("%s", k))
				return nil
			})
			return err
		}
		return nil
	})
	return err
}

// CreateBucket creates a bucket if it doesn't exist.
func (s Storage) CreateBucket(bucket string) error {
	if !s.Opened {
		return fmt.Errorf("db must be opened before creating bucket")
	}
	err := s.DB.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return fmt.Errorf("Error creating bucket : %s", err)
		}
		return nil
	})
	return err
}

// MonitStorage is the general storage associated to the bot.
// It should be available to any plugin, middleware or any other part of the program.
var MonitStorage = Storage{}
