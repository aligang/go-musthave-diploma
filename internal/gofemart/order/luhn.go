package order

func CheckLuhn(number uint64) bool {
	var luhn uint64

	for i := 0; number > 0; i++ {
		cur := number % 10

		if i%2 == 0 { // even
			cur = cur * 2
			if cur > 9 {
				cur = cur%10 + cur/10
			}
		}
		luhn += cur
		number = number / 10
	}
	return luhn%10 == 0
}
