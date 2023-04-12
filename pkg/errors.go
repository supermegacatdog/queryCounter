package pkg

import "errors"

var (
	NotEnoughArgs   = errors.New("not enough args")
	StringEmpty     = errors.New("string is empty")
	StringIncorrect = errors.New("string is empty")

	InvalidMode = errors.New("mode is invalid")
)
