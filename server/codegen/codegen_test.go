package codegen

import (
	"testing"
	"fmt"
)

func Test_defaultChars(_ *testing.T) {
	fmt.Println(string(createAZ09Slice()))
}

func Test_alternateBuilder(_ *testing.T) {
	allChars := []rune{'a', 'b', 'c', 'd'}
	codes := genCodes(allChars, 3)
	for _, code := range codes {
		fmt.Println(string(code))
	}
	fmt.Println("num codes:", len(codes))
}