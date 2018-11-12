package auth

import (
	"database/sql"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"nick-anderssohn-website/ops-gateway/serverutil"
	"os"
	"strings"
	"time"
)

var KibanaEndpoint = &serverutil.Endpoint{
	Path: "/Kibana",
	HandleFunc: handleKibana,
}

var kibanaUrl *url.URL

var loginPageUrl string

func init() {
	log.Println("init kibana url")
	kibanaUrl, _ = url.Parse("http://localhost:5601/app/kibana")
	log.Println("init kibana reverse proxy")
	kibanaReverseProxy = httputil.NewSingleHostReverseProxy(kibanaUrl)
	kibanaReverseProxy.Director = func(req *http.Request) {
		req.URL.Scheme = kibanaUrl.Scheme
		req.URL.Host = kibanaUrl.Host
		req.URL.Path = strings.Replace(singleJoiningSlash(kibanaUrl.Path, req.URL.Path), "/ops/Kibana", "", 1)
		if kibanaUrl.RawQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = kibanaUrl.RawQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = kibanaUrl.RawQuery + "&" + req.URL.RawQuery
		}
		if _, ok := req.Header["User-Agent"]; !ok {
			// explicitly disable User-Agent so it's not set to default value
			req.Header.Set("User-Agent", "")
		}
	}
	loginPageUrl = os.Getenv("LOGIN_PAGE_URL")
}

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}

var kibanaReverseProxy *httputil.ReverseProxy

// First perform authn/authz then pass requests along to kibana
func handleKibana(respWriter http.ResponseWriter, req *http.Request) {
	defer recoverAndLog()

	log.Println("in handle kibana")

	// Todo: send to a different page if login is successful, but they do not have permission.
	user, isAuthorized := handleAuthn(req)
	if !isAuthorized {
		log.Println(loginPageUrl)
		http.Redirect(respWriter, req, loginPageUrl, http.StatusSeeOther)
		return
	}

	log.Println("through authn")
	log.Println("user:", user)

	if !user.CanViewKibana {
		return
	}

	kibanaReverseProxy.ServeHTTP(respWriter, req)
}

// Returns true if okay to pass through to dest.
func handleAuthn(req *http.Request) (*userEntity, bool) {
	sessionCookie, err := req.Cookie(sessionCookieName)
	log.Println("session cookie:", sessionCookie)

	// Check if we were able to get the cookie
	if err != nil {
		if err == http.ErrNoCookie {
			log.Println("not logged in")
		} else {
			log.Println("unknown error ", err)
		}
		return nil, false
	}

	// Check if expired
	if time.Now().Before(sessionCookie.Expires) {
		log.Println("session expired")
		return nil, false
	}

	// Get session
	session, err := grabSessionFromId(sessionCookie.Value)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("session not found")
		} else {
			log.Println("unknown error: ", err)
		}
		return nil, false
	}

	// Get user
	user, err := session.getUser()
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("user not found")
		} else {
			log.Println("unknown error ", err)
		}
		return nil, false
	}

	return user, true
}