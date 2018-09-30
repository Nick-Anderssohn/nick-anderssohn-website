package api

import (
	"bufio"
	"fmt"
	"os"
	"testing"
)

func Test_makeUrl(t *testing.T) {
	w := bufio.NewWriterSize(os.Stdout, 1000)
	fmt.Fprint(w, "Hello, ")
	fmt.Fprint(w, "world!!!")
	w.Flush() // Don't forget to flush!
}
