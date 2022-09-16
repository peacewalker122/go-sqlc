package util

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword (pass string) (string,error){
	hash,err := bcrypt.GenerateFromPassword([]byte(pass),bcrypt.DefaultCost)
	if err != nil {
		return "",fmt.Errorf("can't hash pass, due of: %v", err)
	}
	return string(hash),nil
}

func CheckPass (pass,hashed string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed),[]byte(pass))
}