package codegen

const defaultCodeLen = 4

type Generator struct {
	availableCodes []string
	codeLen int
	availableCharacters []rune
}

func NewDefaultGenerator() *Generator {
	g := &Generator{
		availableCharacters: createAZaz09Slice(),
		codeLen: defaultCodeLen,
	}
	g.instantiateCodes()
	return g
}

func (g *Generator) instantiateCodes() {
	codeMap := map[string]bool{}
	curCombo := make([]rune, g.codeLen)
	permuteAllCombinations(g.availableCharacters, curCombo, 0, g.codeLen, codeMap)
	g.availableCodes = getKeys(codeMap)
}

func createAZaz09Slice() (chars []rune) {
	chars = appendRangeOfChars(chars, 'A', 'Z')
	chars = appendRangeOfChars(chars, 'a', 'z')
	chars = appendRangeOfChars(chars, '0', '9')

	return
}

func appendRangeOfChars(chars []rune, startChar, endChar rune) []rune{
	for c := startChar; c <= endChar; c++ {
		chars = append(chars, c)
	}
	return chars
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