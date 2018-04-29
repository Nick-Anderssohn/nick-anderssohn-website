package codegen

import (
	"fmt"
	"testing"
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

func Test_checkDuplicates(t *testing.T) {
	ng := NewDefaultGenerator()

	checker := map[string]bool{}
	var duplicateCount int
	fmt.Println("len of available codes: ", len(ng.availableCodes))
	for _, code := range ng.availableCodes {
		if checker[code] {
			duplicateCount++
		}
		checker[code] = true
	}
	fmt.Println("num distinct: ", len(checker))
	if duplicateCount > 0 {
		t.Error("There were ", duplicateCount, "duplicates")
	}
}
