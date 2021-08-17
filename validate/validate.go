package validate

const PasswordMinLength = 8

func IsStringEmpty(text string) bool {
	return text == ""
}

func IsNumberNegative(number int) bool {
	return number < 0
}

func IsAge12AndAbove(number int) bool {
	return number >= 12
}

func IsPasswordLengthCorrect(password string) bool {
	return len(password) >= PasswordMinLength
}
