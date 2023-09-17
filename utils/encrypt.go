package utils

import "golang.org/x/crypto/bcrypt"

func EncodeBcrypt(str string) (hs string, err error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(str), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func CompareBcrypt(hash, str string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(str))
	return err == nil
}
