package goerror

import "errors"

var (
	UserNotFound        = errors.New("user not found")
	UsernameHasTaken    = errors.New("username has taken")
	PhoneNumberHasTaken = errors.New("phone number has taken")
	EmailHasTaken       = errors.New("email has taken")

	FeedNotFound = errors.New("feed not found")

	AuthEitherPhoneNumberEmail = errors.New("choose one between phone number and email")

	DataNotFound = errors.New("data not found")
)
