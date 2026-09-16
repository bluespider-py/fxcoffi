package main

import (
	"log"
	"fmt"
	bolt "go.etcd.io/bbolt"
)


func main() {
	// Open the my.db data file in your current directory.
	// It will be created if it doesn't exist.
	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucket([]byte("rootBucket"))
		if err != nil {
			return fmt.Errorf("Bucket creation failed with error - %s", err)
		}

		bucket.Put([]byte("answer"), []byte("42"))
		if err != nil {
			return err
		}
		return nil
	})

	db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("rootBucket"))
		if err != nil {
			return fmt.Errorf("Bucket fetching failed with error - %s", err)
		}

		value := bucket.Get([]byte("answer"))
		fmt.Println("The value is - ", string(value))
		return nil
	})

	
}
