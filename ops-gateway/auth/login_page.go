package auth

import (
	"net/http"
	"nick-anderssohn-website/ops-gateway/serverutil"
)

const LoginPagePath = "/Login/"

var loginFileServerHandler = http.StripPrefix(LoginPagePath, http.FileServer(http.Dir("/resources")))

var LoginPageEndpoint = &serverutil.Endpoint{
	Path: LoginPagePath,
	HandleFunc: loginFileServerHandler.ServeHTTP,
}