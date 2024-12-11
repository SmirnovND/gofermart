package luna

import "strconv"

func LunaAlgorithm(orderNumber string) bool {
	var digits []int
	for _, ch := range orderNumber {
		digit, err := strconv.Atoi(string(ch))
		if err != nil {
			return false // если символ не цифра, возвращаем false
		}
		digits = append(digits, digit)
	}

	var sum int
	for i := len(digits) - 1; i >= 0; i-- {
		digit := digits[i]

		if (len(digits)-i)%2 == 0 {
			digit *= 2
			if digit > 9 {
				digit = digit - 9
			}
		}
		sum += digit
	}

	return sum%10 == 0
}
