package yBytes

/* Основной интерфейс */

func GetCountOnBytes64(rule uint64) int {
	if rule == 0 {
		return 0
	}

	out := 0
	mask := uint64(0xffffffff)
	lowBytes := rule & mask
	upBytes := (rule >> 32) & mask

	if lowBytes > 0 {
		out += GetCountOnBytes32(uint32(lowBytes))
	}
	if upBytes > 0 {
		out += GetCountOnBytes32(uint32(upBytes))
	}
	return out
}

func GetCountOnBytes32(rule uint32) int {
	if rule == 0 {
		return 0
	}

	out := 0
	lowBytes := rule & 0xffff
	upBytes := (rule >> 16) & 0xffff

	if lowBytes > 0 {
		out += GetCountOnBytes16(uint16(lowBytes))
	}
	if upBytes > 0 {
		out += GetCountOnBytes16(uint16(upBytes))
	}
	return out
}

func GetCountOnBytes16(rule uint16) int {
	if rule == 0 {
		return 0
	}
	out := 0
	lowByte := rule & 0xff
	upByte := (rule >> 8) & 0xff

	if lowByte > 0 {
		out += GetCountOnByte8(byte(lowByte))
	}

	if upByte > 0 {
		out += GetCountOnByte8(byte(upByte))
	}

	return out
}

func GetCountOnByte8(b byte) int {
	if b == 0 {
		return 0
	}

	out := 0
	for i := 0; i < 8; i++ {
		mask := 1 << i
		if byte(mask)&b > 0 {
			out++
		}
	}
	return out
}

// GetPositionsOnBytesUI - находим позиции включенных битов в числовых без знаковых линейках
func GetPositionsOnBytesUI(rule interface{}) []int {
	var i, length uint8
	var tmpRule uint64
	var pos int

	switch rule.(type) {
	case uint64:
		length = 8
		tmpRule = rule.(uint64)
	case uint32:
		length = 4
		tmp := rule.(uint32)
		tmpRule = uint64(tmp)
	case uint16:
		length = 2
		tmp := rule.(uint16)
		tmpRule = uint64(tmp)
	case uint8:
		length = 1
		tmp := rule.(uint8)
		tmpRule = uint64(tmp)
	default:
		// не верный тип
		return nil
	}
	if tmpRule == 0 {
		// пустая линейка
		return nil
	}
	out := make([]int, 0, 2)

	// цикл по байтам
	for i = 0; i < length; i++ {
		shift := i * 8
		part := (tmpRule >> shift) & 0xff
		if part == 0 {
			continue
		}
		// цикл по битам
		for bitI := 0; bitI < 8; bitI++ {
			mask := 1 << bitI
			if part&uint64(mask) > 0 {
				pos = (int(i) * 8) + bitI
				out = append(out, pos)
			}
		}
	}

	return out
}
