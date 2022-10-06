package gapi

import (
	"fmt"
	"net/mail"
	"regexp"
)

var (
	Usernamecheck = regexp.MustCompile(`^[a-z0-9_]+$`).MatchString
	Fullnamecheck = regexp.MustCompile(`^[a-zA-Z\\s]+$`).MatchString
)

func validateString(target string, minChar, maxChar int) error {
	if len(target) < minChar || len(target) > maxChar {
		return fmt.Errorf("invalid length of string, must contain %d-%d character", minChar, maxChar)
	}
	return nil
}

func validateUsername(username string) error {
	if err := validateString(username, 3, 100); err != nil {
		return err
	}
	if !Usernamecheck(username) {
		return fmt.Errorf("illegal symbol can't proceed, must contain letter,number,underscore")
	}
	return nil
}

func validatePassword(pass string) error {
	if err := validateString(pass, 5, 100); err != nil {
		return err
	}
	return nil
}

func validateEmail(email string) error {
	if err := validateString(email, 5, 100); err != nil {
		return err
	}
	_, err := mail.ParseAddress(email)
	if err != nil {
		return fmt.Errorf("not valid email: %v", err)
	}
	return nil
}

func validateFullname(fullname string) error {
	if err := validateString(fullname, 3, 100); err != nil {
		return err
	}
	if !Fullnamecheck(fullname) {
		return fmt.Errorf("illegal symbol can't proceed, must contain letter and space")
	}
	return nil
}
