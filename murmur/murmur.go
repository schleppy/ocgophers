package utils

func Hash64WithSeed(data []byte, seed int) int64 {
	length := uint64(len(data))
	var m uint64 = 0xc6a4a7935bd1e995
	var r uint64 = 47

	h := uint64(seed&0xffffffff) ^ (length * m)

	length8 := length / 8
	var i uint64
	for i = 0; i < length8; i++ {
		i8 := i * 8
		k := uint64(int64(data[i8+0]&0xff) +
			int64(data[i8+1]&0xff)<<8 +
			int64(data[i8+2]&0xff)<<16 +
			int64(data[i8+3]&0xff)<<24 +
			int64(data[i8+4]&0xff)<<32 +
			int64(data[i8+5]&0xff)<<40 +
			int64(data[i8+6]&0xff)<<48 +
			int64(data[i8+7]&0xff)<<56)
		k *= m
		k ^= k >> r
		k *= m

		h ^= k
		h *= m
	}

	h2 := int64(h)

	switch length % 8 {
	case 7:
		h2 ^= int64(int64(data[(int64(length)&-8)+6]&0xff) << 48)
		fallthrough
	case 6:
		h2 ^= int64(int64(data[(int64(length)&-8)+5]&0xff) << 40)
		fallthrough
	case 5:
		h2 ^= int64(int64(data[(int64(length)&-8)+4]&0xff) << 32)
		fallthrough
	case 4:
		h2 ^= int64(int64((data[(int64(length)&-8)+3] & 0xff)) << 24)
		fallthrough
	case 3:
		h2 ^= int64(int64(data[(int64(length)&-8)+2]&0xff) << 16)
		fallthrough
	case 2:
		h2 ^= int64(int64(data[(int64(length)&-8)+1]&0xff) << 8)
		fallthrough
	case 1:
		h2 ^= int64(int64(data[int64(length)&-8] & 0xff))
		h2 = int64(uint64(h2) * uint64(m))
	}

	h2 ^= int64(uint64(h2) >> r)
	h2 = int64(uint64(h2) * m)
	h2 ^= int64(uint64(h2) >> r)

	return h2
}

func Hash64(data string) int64 {
	seed := 0xe17a1465
	return Hash64WithSeed([]byte(data), seed)
}
