
package bitcask

import (
	"time"
	"encoding/binary"
)

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
	
	keyLen := binary.BigEndian.Uint32(data[8:12])
	valueLen := binary.BigEndian.Uint32(data[12:16])

	return data[16:16+keyLen], data[16+keyLen:16+keyLen+valueLen], nil
	
}


