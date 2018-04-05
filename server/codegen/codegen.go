package codegen

type Generator struct {
	availableCodes []string
	codeLen int
	availableCharacters []string
}

func getKeys(m map[string]bool) []string {
	keys := make([]string, 0, len(m))

	for k := range m {
		keys = append(keys, k)
	}

	return keys
}

// len of curCombo must equal comboSize
func permuteAllCombinations(allChars, curCombo []rune, index, comboSize int, permutationMap map[string]bool) {
	if len(curCombo) != comboSize {
		return
	}

	// if a combo is reached, get all permutations for that combo
	if index == comboSize {
		permute(permutationMap, string(curCombo), 0, len(curCombo) - 1)
		return
	}

	for i := 0; i <= len(allChars)-1 && len(allChars)-i >= comboSize-index; i++ {
		curCombo[index] = allChars[i]
		permuteAllCombinations(allChars[1:], curCombo, index+1, comboSize, permutationMap)
	}
}

// Does NOT store duplicates
func permute(permutationMap map[string]bool, curPermutation string, startIndex, endIndex int) {
	if startIndex == endIndex {
			permutationMap[curPermutation] = true
	} else {
		for i := startIndex; i <= endIndex; i++ {
			curPermutation = swap(curPermutation, i, startIndex)
			permute(permutationMap, curPermutation, startIndex+1, endIndex)
			curPermutation = swap(curPermutation, i, startIndex)
		}
	}
}

func swap(str string, i1, i2 int) string {
	chars := []rune(str)
	chars[i1], chars[i2] = chars[i2], chars[i1]
	return string(chars)
}