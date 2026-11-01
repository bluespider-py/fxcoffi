
package bitcask

import (
	"time"
	"encoding/binary"
	"errors"
	"hash/crc32"
)

const headerSize = 20  // checksum(4) + timestamp(8) + keySize(4) + valueSize(4)
const keyLenIndex = 12
const valueLenIndex = 16
const dataInitIndex = 20

var ErrTruncatedRecord = errors.New("record is truncated or has inconsistent sizes")
var ErrChecksumMismatch = errors.New("checksum of the record failed")

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

	body := make([]byte, 0, 4 + len(timeByte) + len(keyLenByte) + len(valueLenByte) + len([]byte(key)) + len([]byte(value)))
	body = append(body, timeByte...)
	body = append(body, keyLenByte...)
	body = append(body, valueLenByte...)
	body = append(body, []byte(key)...)
	body = append(body, []byte(value)...)

	checksum := uint32ToBytes(crc32.ChecksumIEEE(body))
	body = append(checksum, body...)
	
	return body
}

func decodeRecord(data []byte) (key, value []byte, err error) {

	// 0 -> 3 is checksum
	// 4 -> 11 is timestamp
	// 12 - 15 is key len
	// 16 - 19 is value len
	// 20 onwards is key and value

	if len(data) < headerSize {
		return nil, nil, ErrTruncatedRecord
	}
	
	keyLen := uint64(binary.BigEndian.Uint32(data[keyLenIndex:keyLenIndex+4]))
	valueLen := uint64(binary.BigEndian.Uint32(data[valueLenIndex:valueLenIndex+4]))

	if uint64(len(data)) < uint64(headerSize)+keyLen+valueLen {
		return nil, nil, ErrTruncatedRecord
	}

	// checksum/data corruption validation
	if binary.BigEndian.Uint32(data[0:4]) != crc32.ChecksumIEEE(data[4:]) {
		return nil, nil, ErrChecksumMismatch
	}	

	return data[headerSize : headerSize+keyLen],
		data[headerSize+keyLen : headerSize+keyLen+valueLen],
		nil	
}


