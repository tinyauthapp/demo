package main

import (
	_ "embed"
	"html/template"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
)

//go:embed page.html
var pageHTML string

//go:embed background.jpg
var background []byte

var tmpl = template.Must(template.New("page").Parse(pageHTML))

// identityHeaders are the headers Tinyauth sets on authenticated requests.
var identityHeaders = []string{"Remote-User", "Remote-Name", "Remote-Email", "Remote-Groups"}

type row struct{ Key, Value string }

type page struct {
	Protected bool     // false on the public landing page, true on /protected
	Identity  []row    // the four Tinyauth headers, in display order
	Groups    []string // Remote-Groups split into individual badges
	Headers   []row    // every request header, sorted by name
	LogoutURL string
}

// render builds the view model from the request headers and writes the page.
func render(w http.ResponseWriter, r *http.Request, protected bool, logoutURL string) {
	p := page{Protected: protected, LogoutURL: logoutURL}
	for _, h := range identityHeaders {
		p.Identity = append(p.Identity, row{h, r.Header.Get(h)})
	}
	for g := range strings.SplitSeq(r.Header.Get("Remote-Groups"), ",") {
		if g = strings.TrimSpace(g); g != "" {
			p.Groups = append(p.Groups, g)
		}
	}

	keys := make([]string, 0, len(r.Header))
	for k := range r.Header {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		p.Headers = append(p.Headers, row{k, strings.Join(r.Header[k], ", ")})
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store") // page reflects per-request headers
	// w.Header().Set("X-Content-Type-Options", "nosniff")
	if err := tmpl.Execute(w, p); err != nil {
		log.Printf("render: %v", err)
	}
}

func backgroundHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Cache-Control", "public, max-age=3600")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Write(background)
}

func main() {
	logoutURL := os.Getenv("LOGOUT_URL")
	if logoutURL == "" {
		logoutURL = "/logout"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) { render(w, r, false, logoutURL) })
	mux.HandleFunc("GET /protected", func(w http.ResponseWriter, r *http.Request) { render(w, r, true, logoutURL) })
	mux.HandleFunc("GET /background.jpg", backgroundHandler)

	addr := ":3000"
	log.Printf("tinyauth-demo listening on %s (logout -> %s)", addr, logoutURL)
	log.Fatal(http.ListenAndServe(addr, mux))
}
