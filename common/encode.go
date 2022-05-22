package common

import (
	"crypto/sha256"
	"fmt"
)

//GetHashPassword 将密码sha256加密
func GetHashPassword(password, salt string) string {
	data := []byte(password + salt)
	hashPwd := sha256.Sum256(data)
	res := fmt.Sprintf("%x", hashPwd)
	return res

}
