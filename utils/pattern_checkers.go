package utils

import "regexp"

const uuidPattern = "^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$"

const emailPattern = "^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9]" +
	"(?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"

//TODO: Format all phone numbers to one structure
const mobilePhonePattern = `^(?:(?:\(?(?:00|\+)([1-4]\d\d|[1-9]\d?)\)?)?[\-\.\ \\\/]?)?((?:\(?\d{1,}\)?[\-\.\ \\\/]?){0,})(?:[\-\.\ \\\/]?(?:#|ext\.?|extension|x)[\-\.\ \\\/]?(\d+))?$`

func IsUUID(UUID string) bool {
	res, _ := regexp.MatchString(uuidPattern, UUID)
	return res
}

func IsEmail(email string) bool {
	res, _ := regexp.MatchString(emailPattern, email)
	return res
}

func IsPhoneNumber(number string) bool {
	res, _ := regexp.MatchString(mobilePhonePattern, number)
	return res
}
