package yBytes

// хэш слова, первые 4 буквы или последние 4 буквы
type hashWord uint32

// хэш слова с доп информацией о длине слова, и вариациях расположения ударения в одинаково звучащих словах

func packH3(h3 *uint32, code uint32) {
	var tmpH3, tmpHeader uint32

	if (*h3) == 0 {
		*h3 = code << 24
		return
	}
	tmpH3 = (*h3) >> 8
	tmpHeader = code << 24
	tmpActiveH3 := tmpH3 & 0xffffff00
	*h3 = tmpHeader | tmpActiveH3
}

// packSpecialH3 - добавить доп информацию в третий хэш
// параметры hp, tp, позволяют скомбинировать 3 варианта записи слов которые одинаково пишутся но имеют разное смысловое
// значение
// hp - header part - ударение приходится на первую половину слова
// tp - tail part - ударение приходится на вторую половину слова
// возврат:
// 	- 0 - все хорошо
//  - 1 - длина слова превышает допустимый лимит 63 символа
//  - 2 - hp и ht были установлены в true. Предпочтение было отдано hp
func packSpecialH3(h3 *uint32, lenWord uint32, hp, tp bool) int {
	if hp && tp {
		return 2
	}
	// больше 63 символов
	if lenWord > 0x3f {
		return 1
	}
	if hp {
		*h3 |= 1 << 7
	}
	if tp {
		*h3 |= 1 << 6
	}
	*h3 |= lenWord
	return 0
}
