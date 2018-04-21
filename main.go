package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/facebook"
	"github.com/stretchr/gomniauth/providers/github"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/objx"

	"gitlab.com/Orenge/chat/trace"
)

type templateHandler struct {
	// Web servers in Go are automatically concurrent and once our chat application takes
	// the world by storm, we could very well expect to have many concurrent calls to the
	// ServeHTTP method.
	once     sync.Once //compile the template once regardless of how many goroutines are calling ServeHTTP
	filename string
	templ    *template.Template //keep the reference to the compiled template
}

func (t templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

func main() {

	var addr = flag.String("addr", ":8000", "Listen at what port?")
	flag.Parse()
	gomniauth.SetSecurityKey("ThisMySaltp4$$~0d")
	gomniauth.WithProviders(
		facebook.New("key", "secret",
			"http://localhost:3000/auth/callback/facebook"),
		github.New("key", "secret",
			"http://localhost:3000/auth/callback/github"),
		google.New("", "",
			"http://localhost:3000/auth/callback/google"),
	)
	r := newRoom(UseGravatar)
	r.tracer = trace.New(os.Stdout)
	//go the room going
	go r.run()
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"}) //MustAuth method for our login page because
	//it will cause an infinite redirection loop
	http.HandleFunc("/auth/", loginHandler)
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

	http.Handle("/upload", &templateHandler{filename: "upload.html"})
	http.HandleFunc("/uploader", uploaderHandler)

	http.Handle("/room", r)

	// FileServer - will simply serve static files, provide index listings, and
	// generate the 404 Not Found error if it cannot find the file. The http.Dir function
	// allows us to specify which folder we want to expose publicly
	http.Handle("/avatars/", http.StripPrefix("/avatars/", http.FileServer(http.Dir("./avatars"))))

	log.Println("Starting to serve at", *addr)

	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServeErr:", err)
	}
}
