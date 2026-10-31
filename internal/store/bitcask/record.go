
package bitcask

import (
	"time"
	"encoding/binary"
	"errors"
)

const headerSize = 16 // timestamp(8) + keySize(4) + valueSize(4); becomes 20 when you add the CRC
var ErrTruncatedRecord = errors.New("record is truncated or has inconsistent sizes.")
var ErrChecksumMismatch = errors.New("checksum of the record failed.")

func currentTimeUnix() uint64 {
	return uint64(time.Now().Unix())
}

func uint64ToBytes(el uint64) []byte {
	bytes := make([]byte, 8)
	binary.BigEndian.PutUint64(bytes, el)
	return bytes 
}

func uint32ToBytes(el uint32) []byte {
	bytes := make([]byte, 4)
	binary.BigEndian.PutUint32(bytes, el)
	return bytes
}

func encodeRecord(key, value []byte) []byte {
	unixTime := currentTimeUnix()
	timeByte := uint64ToBytes(unixTime)
	keyLenByte := uint32ToBytes(uint32(len(key)))
	valueLenByte := uint32ToBytes(uint32(len(value)))

	body := make([]byte, 0, len(timeByte) + len(keyLenByte) + len(valueLenByte) + len([]byte(key)) + len([]byte(value)))
	body = append(body, timeByte...)
	body = append(body, keyLenByte...)
	body = append(body, valueLenByte...)
	body = append(body, []byte(key)...)
	body = append(body, []byte(value)...)
	
	return body

}

func decodeRecord(data []byte) (key, value []byte, err error) {
	
	// 0 -> 7 is timestamp
	// 8 - 11 is key len
	// 12 - 15 is value len
	// 16 onwards is key and value

	if len(data) < headerSize {
		return nil, nil, ErrTruncatedRecord
	}

	
	keyLen := uint64(binary.BigEndian.Uint32(data[8:12]))
	valueLen := uint64(binary.BigEndian.Uint32(data[12:16]))

	if uint64(len(data)) < uint64(headerSize)+keyLen+valueLen {
		return nil, nil, ErrTruncatedRecord
	}

	return data[headerSize : headerSize+keyLen],
		data[headerSize+keyLen : headerSize+keyLen+valueLen],
		nil	
}


