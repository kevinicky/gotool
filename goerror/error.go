package goerror

import "errors"

var (
	UserNotFound        = errors.New("user not found")
	UsernameHasTaken    = errors.New("username has taken")
	PhoneNumberHasTaken = errors.New("phone number has taken")

	FeedNotFound = errors.New("feed not found")

	DataNotFound = errors.New("data not found")
)
