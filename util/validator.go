package util

import "regexp"

func EmailValid(email string) (bool,error) {
	return regexp.MatchString("^([A-Za-z0-9])+@([A-Za-z0-9])+\\.([A-Za-z]{2,4})$",email)
}
func PhoneValid(phone string) (bool,error) {
	return regexp.MatchString("^1[34578]\\d{9}$",phone)
}
func UsernameValid(username string) (bool,error) {
	return regexp.MatchString("^[a-zA-Z0-9\u4E00-\u9FFF]{1,20}$",username)
}