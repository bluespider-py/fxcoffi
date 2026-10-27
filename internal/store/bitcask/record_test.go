package bitcask

import (
	"fmt"
	"testing"
	"errors"
)

func TestBasicRecord(t *testing.T) {
	key := []byte("answer")
	value := []byte("cat")

	encoded := encodeRecord(key, value)
	gotKey, gotValue, err := decodeRecord(encoded)

	fmt.Printf("key - %q, value - %q, error - %v\n", gotKey, gotValue, err)

	if err != nil {
		t.Fatalf("Function threw error - %v\n", err)
	}

	if string(key) != string(gotKey) {
		t.Errorf("Keys mismatch - %v against %v\n", key, gotKey)
	}

	if string(value) != string(gotValue) {
		t.Errorf("Values mismatch - %v against %v\n", value, gotValue)
	}
}

func TestCorruptedEncoding(t *testing.T) {

	key  := []byte("answer") 
	var value []byte

	encoded := encodeRecord(key, value)
	encoded[len(encoded)/2] ^= 0xFF
	
	_, _, err := decodeRecord(encoded)

	if err == nil {
		t.Fatal("expected an error for a corrupted record, got nil")
	}
	if !errors.Is(err, ErrTruncatedRecord) {
		t.Errorf("expected ErrTruncatedRecord, got %v", err)
	}

}

	

