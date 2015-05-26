package tonic

import "log"

func lastChar(str string) uint8 {
	size := len(str)
	if size == 0 {
		log.Panic("The length of the string can't be 0")
	}
	return str[size-1]
}
