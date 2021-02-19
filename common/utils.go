package common

import "encoding/binary"

const TIME_FORMAT = "2006-01-02 15:04:05"

func GenIntFromType(data []byte) uint64 {
	return uint64(binary.BigEndian.Uint32(data[:]))
}

func GenIntFromLength(data []byte) uint64 {
	return uint64(binary.BigEndian.Uint64(data[:]))
}

func GenLengthFromInt(l int) [8]byte {
	var result [8]byte
	binary.BigEndian.PutUint64(result[:], uint64(l))
	return result
}

func GenTypeFromInt(t int) [4]byte {
	var result [4]byte
	binary.BigEndian.PutUint32(result[:], uint32(t))
	return result
}

func StringInSlice(s string, ss []string) bool {
	for idx := range ss {
		if s == ss[idx] {
			return true
		}
	}
	return false
}

func Int64InSlice(i int64, ii []int64) bool {
	for idx := range ii {
		if i == ii[idx] {
			return true
		}
	}
	return false
}

func Int64SliceEqual(a []int64, b []int64) bool {
	if len(a) != len(b) {
		return false
	}

	for idx := range a {
		av := a[idx]
		if Int64InSlice(av, b) == false {
			return false
		}
	}

	return true
}
