package unpacker

import (
	"errors"
	"strconv"
	"unicode"
)

func UnpackString(s string) (string, error) {
	var result []rune
	runes := []rune(s)
	length := len(runes)
	if length == 0 {
		return "", nil
	}
	if unicode.IsDigit(runes[0]) {
		return "", errors.New("shouldn't start with a digit")
	}
	i := 0
	for i < length {
		current := runes[i]
		if current == '\\' {
			if i+1 < length {
				result = append(result, runes[i+1])
				i += 2
				continue
			} else {
				result = append(result, '\\')
				i++
				continue
			}
		}
		if !unicode.IsDigit(current) {
			result = append(result, current)
			i++
			continue
		}
		start := i
		for i < length && unicode.IsDigit(runes[i]) {
			i++
		}
		numStr := string(runes[start:i])
		count, err := strconv.Atoi(numStr)
		if err != nil {
			return "", errors.New("not numeric")
		}
		if count == 0 {
			if len(result) > 0 {
				result = result[:len(result)-1]
			}
		} else if len(result) > 0 {

			lastChar := result[len(result)-1]
			for j := 1; j < count; j++ {
				result = append(result, lastChar)
			}
		}
	}
	return string(result), nil
}