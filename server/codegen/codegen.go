package codegen

import (
	"math"
	"math/rand"

	"fmt"

	"gopkg.in/cheggaaa/pb.v2"
)

const defaultCodeLen = 4

var totalToGen int64
var progressBar *pb.ProgressBar

type Generator struct {
	availableCodes  []string
	codeLen         int
	allChars        []rune
	RequestCodeChan chan bool
	CodeChan        chan string
}

func NewDefaultGenerator() *Generator {
	g := &Generator{
		allChars:        createAZ09Slice(),
		codeLen:         defaultCodeLen,
		RequestCodeChan: make(chan bool),
		CodeChan:        make(chan string),
	}
	g.instantiateCodes()
	return g
}

func (g *Generator) instantiateCodes() {
	totalToGen = int64(math.Pow(float64(len(g.allChars)), float64(g.codeLen)))
	progressBar = pb.StartNew(int(totalToGen))
	g.availableCodes = genCodes(g.allChars, g.codeLen)
	progressBar.Finish()
}

func (g *Generator) GetCode() string {
	i := rand.Intn(len(g.availableCodes))
	fmt.Println("random index: ", i)
	code := g.availableCodes[i]
	g.availableCodes = append(g.availableCodes[:i], g.availableCodes[i+1:]...)
	return code
}

func (g *Generator) RestoreCode(code string) {
	g.availableCodes = append(g.availableCodes, code)
}

func (g *Generator) RunAndWaitForCodeRequests() {
	for {
		<-g.RequestCodeChan
		g.CodeChan <- g.GetCode()
	}
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

func genCodes(allChars []rune, codeLen int) []string {
	var allCodes []string
	genCodesR(allChars, "", codeLen, &allCodes)
	return allCodes
}

func genCodesR(allChars []rune, curCode string, k int, storage *[]string) {
	if k == 0 {
		*storage = append(*storage, curCode)
		if progressBar != nil {
			progressBar.Increment()
		}
		return
	}

	for _, c := range allChars {
		newCode := curCode + string(c)
		genCodesR(allChars, newCode, k-1, storage)
	}
}
