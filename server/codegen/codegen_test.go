package codegen

import (
	"testing"
	"fmt"
)

func Test_permute(_ *testing.T) {
	str := "abc"
	permutationMap := map[string]bool{}
	permute(permutationMap, str, 0, len(str) - 1)
	fmt.Println(getKeys(permutationMap))
}

func Test_permuteAllCombinations(_ *testing.T) {
	comboSize := 2
	allChars := []rune{'a', 'b', 'c', 'd'}
	curCombo := make([]rune, comboSize)
	permutationMap := map[string]bool{}
	permuteAllCombinations(allChars, curCombo, 0, comboSize, permutationMap)
	fmt.Println(getKeys(permutationMap))
}