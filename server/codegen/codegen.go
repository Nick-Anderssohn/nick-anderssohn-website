package codegen

import (
	"math"
	"math/rand"

	"gopkg.in/cheggaaa/pb.v2"
)

const defaultCodeLen = 4

var totalToGen int64
var genCount int64
var progressBar *pb.ProgressBar

type Generator struct {
	availableCodes [][]rune
	codeLen        int
	allChars       []rune
}

func NewDefaultGenerator() *Generator {
	g := &Generator{
		allChars: createAZ09Slice(),
		codeLen:  defaultCodeLen,
	}
	g.instantiateCodes()
	return g
}

func (g *Generator) instantiateCodes() {
	totalToGen = int64(math.Pow(float64(len(g.allChars)), float64(g.codeLen)))
	genCount = 0
	progressBar = pb.StartNew(int(totalToGen))
	g.availableCodes = genCodes(g.allChars, g.codeLen)
	progressBar.Finish()
}

func (g *Generator) GetCode() string {
	i := rand.Intn(len(g.availableCodes))
	code := g.availableCodes[i]
	g.availableCodes = append(g.availableCodes[:i], g.availableCodes[i+1:]...)
	return string(code)
}

func (g *Generator) RestoreCode(code string) {
	g.availableCodes = append(g.availableCodes, []rune(code))
}

func createAZ09Slice() (chars []rune) {
	chars = appendRangeOfChars(chars, 'A', 'Z')
	chars = appendRangeOfChars(chars, '0', '9')

	return
}

func appendRangeOfChars(chars []rune, startChar, endChar rune) []rune {
	for c := startChar; c <= endChar; c++ {
		chars = append(chars, c)
	}
	return chars
}

func genCodes(allChars []rune, codeLen int) [][]rune {
	allCodes := make([][]rune, int(math.Pow(float64(len(allChars)), float64(codeLen))))
	for i := 0; i < len(allCodes); i++ {
		allCodes[i] = make([]rune, codeLen)
	}
	for i := 0; i < codeLen; i++ {
		for codeIndex := 0; codeIndex < len(allCodes); codeIndex++ {
			allCodes[codeIndex][i] = allChars[(codeIndex/(i+1))%len(allChars)]
			genCount++
			if genCount%4 == 0 {
				progressBar.Increment()
			}
		}
	}
	return allCodes
}
