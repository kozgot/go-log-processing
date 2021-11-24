package utils

import "log"

// FailOnError fails if the provided error is not nil, and logs the message.
func FailOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
