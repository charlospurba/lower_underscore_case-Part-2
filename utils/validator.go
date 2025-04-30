package utils

import "regexp"

const (
	UsernameMinLength  = 3
	UsernameMaxLength  = 20
	PasswordMinLength  = 8
	NameMinLength      = 3
	NameMaxLength      = 20
	AgeMin             = 15
	EmailGmailSuffix   = "@gmail.com"
	EmailRequired      = true
)

var (
	UsernameRegex = regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	EmailRegex    = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
)

