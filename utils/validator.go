package utils

import "regexp"

// Username validation rules
const (
	UsernameMinLength = 3
	UsernameMaxLength = 20
)

var UsernameRegex = regexp.MustCompile(`^[a-zA-Z0-9]+$`)

// Email validation rules
const (
	EmailRequired    = true
	EmailGmailSuffix = "@gmail.com"
)

var EmailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// Password validation rules
const PasswordMinLength = 8

// Name (FirstName, LastName) validation rules
const (
	NameMinLength = 3
	NameMaxLength = 20
)

// Age validation rules
const AgeMin = 15
