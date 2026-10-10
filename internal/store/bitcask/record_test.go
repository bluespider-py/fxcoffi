package bitcask

import (
	"fmt"
	"testing"
)

func TestDummy(t *testing.T) {
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
