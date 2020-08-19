package global

import (
	"errors"
	"regexp"
)

func VerifyMobileFormat(mobileNum string) (err error) {
	regular := "^1[3-9]\\d{9}$"
	reg := regexp.MustCompile(regular)
	Bool := reg.MatchString(mobileNum)
	if !Bool {
		err = errors.New("手机号格式不正确")
		return
	}
	return
}

func IsValidCitizenNo18(citizenNo18 string) (err error) {
	regular := "^[1-9]\\d{5}(18|19|([23]\\d))\\d{2}((0[1-9])|(10|11|12))(([0-2][1-9])|10|20|30|31)\\d{3}[0-9Xx]$"
	reg := regexp.MustCompile(regular)
	Bool := reg.MatchString(citizenNo18)
	if !Bool {
		err = errors.New("身份证格式不正确")
		return
	}
	return
}
