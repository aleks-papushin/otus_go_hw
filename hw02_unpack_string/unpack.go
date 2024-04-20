package hw02unpackstring

import (
	"errors"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(s string) (string, error) {
	in := []rune(s)
	out := []rune{}

	if len(s) == 0 {
		return s, nil
	}

	if isDigit(in[0]) {
		return "", ErrInvalidString
	}

	for i := 0; i < len(in); i++ {
		e := in[i]

		if !isDigit(e) {
			out = append(out, e)
		} else {
			// вторая цифра подряд
			if isDigit(in[i-1]) {
				return "", ErrInvalidString
			}
			// если это 0, то перезапишем последнюю букву в возвращаемом массиве
			eAsDigit := int(e - '0')
			if eAsDigit == 0 {
				out = out[:len(out)-1]
			} else {
				// - если это цифра более 0, то во вложенном цикле (количество итераций e-1)
				// - записываем в массив такое же значение, как и последнее записанное
				for i := 0; i < eAsDigit-1; i++ {
					out = append(out, out[len(out)-1])
				}
			}
		}
	}

	return string(out), nil
}

func isDigit(r rune) bool {
	return int(r) > 47 && int(r) < 58
}
