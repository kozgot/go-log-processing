package utils

import "log"

// FailOnError checks if the provided error is not nil, and logs the error.
func FailOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
