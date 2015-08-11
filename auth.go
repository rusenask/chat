package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/stretchr/gomniauth"
	"github.com/stretchr/objx"
)

// giving specific name through which this package will be accessed
import gomniauthcommon "github.com/stretchr/gomniauth/common"

// ChatUser exposes information that is needed in order for our Avatar
// implementations to generate the correct URLs.
type ChatUser interface {
	UniqueID() string
	AvatarURL() string
}

// chatUser defines actual implementation that implements the ChatUser interface
// also uses Go's feature - type embedding. This way our struct implements
// the interface automatically
type chatUser struct {
	gomniauthcommon.User
	uniqueID string
}

func (u chatUser) UniqueID() string {
	return u.uniqueID
}

type authHandler struct {
	next http.Handler
}

func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if cookie, err := r.Cookie("auth"); err == http.ErrNoCookie || cookie.Value == "" {
		// user is not authenticated
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusTemporaryRedirect)
	} else if err != nil {
		// some other error
		panic(err.Error())
	} else {
		// success, proceed to next Handler
		h.next.ServeHTTP(w, r)
	}
}

// MustAuth wraps standard handler and enforces authentication
func MustAuth(handler http.Handler) http.Handler {
	return &authHandler{next: handler}
}

// loginHandler handles the third-party login process.
// format: /auth/{action}/{provider}
func loginHandler(w http.ResponseWriter, r *http.Request) {
	segs := strings.Split(r.URL.Path, "/")
	// TODO: handle exception if there are too few segments
	action := segs[2]
	provider := segs[3]
	switch action {
	case "login":
		provider, err := gomniauth.Provider(provider)
		if err != nil {
			log.Fatalln("Eerror when trying to get provider", provider, "-", err)
		}
		loginURL, err := provider.GetBeginAuthURL(nil, nil)
		if err != nil {
			log.Fatalln("error when trying to GetBeginAuthURL for", provider, "-", err)
		}
		// Since http.ResponseWriter is an interface and header is a method on that
		// interface, we have to call it so it will return an http.Header
		// instead of w.Header.Set(...) we use w.Header().Set(...)
		w.Header().Set("Location", loginURL)
		w.WriteHeader(http.StatusTemporaryRedirect)

	case "callback":

		provider, err := gomniauth.Provider(provider)
		if err != nil {
			log.Fatalln("Error when trying to get provider", provider, "-", err)
		}
		// Parsin RawQuery from http.Request into objx.Map (multi-purpose map type
		// that gomniauth uses) and the CompleteAuth method uses the URL query param
		// values to complete the authentication handshake. If all is okay - creds
		// are given
		creds, err := provider.CompleteAuth(objx.MustFromURLQuery(r.URL.RawQuery))
		if err != nil {
			log.Fatalln("Error when trying to complete auth for", provider, "-", err)
		}
		// now accessing user's basic data
		user, err := provider.GetUser(creds)
		if err != nil {
			log.Fatalln("Error when trying to get user from", provider, "-", err)
		}
		fmt.Println("User authenticated: ", user.Name())
		// creating new md5 hasher from crypto package, implements io.Writer interface
		m := md5.New()
		// ensuring email is lower case and generate md5 hash
		// writing a string of bytes to hasher through io.WriteString
		io.WriteString(m, strings.ToLower(user.Name()))
		// caling Sum on hasher returns the current hash for the bytes written
		userID := fmt.Sprintf("%x", m.Sum(nil))
		// encoding user data with Base64 in JSON object
		authCookieValue := objx.New(map[string]interface{}{
			"userid": userID,
			"name":   user.Name(),
			// adding avatar URL to cookie
			"avatar_url": user.AvatarURL(),
			"email":      user.Email(),
		}).MustBase64()
		// storing encoded object in cookie
		http.SetCookie(w, &http.Cookie{
			Name:  "auth",
			Value: authCookieValue,
			Path:  "/"})

		w.Header()["Location"] = []string{"/chat"}
		w.WriteHeader(http.StatusTemporaryRedirect)

	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "auth action %s not supported", action)

	}
}
