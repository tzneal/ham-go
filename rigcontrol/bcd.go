package rigcontrol

func ParseBCD(buf []byte) uint64 {
	var ret uint64
	for _, v := range buf {
		ret *= 10
		hb := (v & 0xF0) >> 4
		ret += uint64(hb)

		ret *= 10
		lb := (v & 0x0f)
		ret += uint64(lb)
	}
	return ret
}

func ToBCD(buf []byte, v uint64, len int) {
	for i := 0; i < len; i++ {
		lb := v % 10
		v /= 10
		hb := v % 10
		v /= 10
		buf[len-i-1] = byte((hb << 4) | lb)
	}
}
