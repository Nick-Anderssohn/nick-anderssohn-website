package api

import (
	"full-share/server/codegen"
	"full-share/server/rest"
	"net/http"
)

const (
	fullShareBase = "/fullShare"
	getCode       = "/getCode"
)

var generator *codegen.Generator

func init() {
	generator = codegen.NewDefaultGenerator()
	go generator.RunAndWaitForCodeRequests()
}

func GetEndpoints() []*rest.Endpoint {
	return []*rest.Endpoint{
		{
			Path:       fullShareBase + getCode,
			HandleFunc: GetCode,
		},
	}
}

func GetCode(writer http.ResponseWriter, _ *http.Request) {
	generator.RequestCodeChan <- true
	code := <-generator.CodeChan
	writer.Write([]byte(code))
}
