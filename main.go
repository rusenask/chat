package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"text/template"
	"trace"

	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/facebook"
	"github.com/stretchr/gomniauth/providers/github"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/objx"
	"github.com/stretchr/signature"
)

// templ represents a single template
type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

// ServeHTTP handles the HTTP Request
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	data := map[string]interface{}{
		"Host": r.Host,
	}
	if authCookie, err := r.Cookie("auth"); err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}
	t.templ.Execute(w, data)
}

// Configuration defining authentication keys for OAuth
// TODO: auth details should be nested under it's own key to make room for more
// configs
type Configuration struct {
	GoogleKey    string
	GoogleSecret string
}

func main() {
	// reading configuration file and creating "configuration" object to store
	// values
	file, _ := os.Open("conf.json")
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("error:", err)
	}
	// looking for option args when starting App
	// like ./chat -addr=":3000" would start on this port
	var addr = flag.String("addr", ":8080", "App address")
	flag.Parse() // parse the flag
	// setting up gomniauth
	// creating random key http://godoc.org/github.com/stretchr/signature#RandomKey
	// in addition to gomniauth package we need to download:
	// go get github.com/clbanning/x2j
	// go get github.com/ugorji/go/codec
	// go get labix.org/v2/mgo/bson
	gomniauth.SetSecurityKey(signature.RandomKey(64))
	gomniauth.WithProviders(
		facebook.New("key", "secret",
			"http://localhost:8080/auth/callback/facebook"),
		github.New("key", "secret",
			"http://localhost:8080/auth/callback/github"),
		google.New(configuration.GoogleKey,
			configuration.GoogleSecret,
			"http://localhost:8080/auth/callback/google"),
	)
	r := newRoom(UseAuthAvatar)
	r.tracer = trace.New(os.Stdout)
	// wrapping /chat handler with MustAuth to enforce authentication
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	// logout, deleting cookie data
	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:   "auth",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		})
		w.Header()["Location"] = []string{"/chat"}
		w.WriteHeader(http.StatusTemporaryRedirect)
	})
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/room", r)
	// get the room going
	go r.run()
	// start the web server
	log.Println("Starting we server on ", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
