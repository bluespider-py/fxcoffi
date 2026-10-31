package bitcask

import (
	"testing"
	"encoding/binary"
	"errors"
)

func TestBasicRecord(t *testing.T) {
	key := []byte("answer")
	value := []byte("cat")

	encoded := encodeRecord(key, value)
	gotKey, gotValue, err := decodeRecord(encoded)

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

func TestDecodeOverflowingSizes(t *testing.T) {
	data := make([]byte, headerSize)
	binary.BigEndian.PutUint32(data[8:12], 0xFFFFFFF0) // keyLen; valueLen stays 0

	_, _, err := decodeRecord(data)

	if !errors.Is(err, ErrTruncatedRecord) {
		t.Errorf("expected ErrTruncatedRecord, got %v", err)
	}
}

func TestCorruptedPayload(t *testing.T) {
	encoded := encodeRecord([]byte("answer"), []byte("cat"))
	encoded[headerSize] ^= 0xFF

	_, _, err := decodeRecord(encoded)

	if !errors.Is(err, ErrChecksumMismatch) {
		t.Errorf("expected ErrChecksumMismatch, got %v", err)
	}
}
