package api

import (
	"full-share/server/codegen"
	"net/http"
	"full-share/server/rest"
)

const (
	fullShareBase = "/fullShare"
	getCode = "/getCode"
)

var generator *codegen.Generator

func init() {
	generator = codegen.NewDefaultGenerator()
}

func GetEndpoints() []*rest.Endpoint {
	return []*rest.Endpoint{
		{
			Path: fullShareBase + getCode,
			HandleFunc: GetCode,
		},
	}
}

func GetCode(writer http.ResponseWriter, _ *http.Request) {
	code := generator.GetCode()
	writer.Write([]byte(code))
}