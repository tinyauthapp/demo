package main

import (
	_ "embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
)

//go:embed page.html
var pageHTML string

//go:embed bg.webp
var bg []byte

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
	w.Header().Set("Cache-Control", "no-store")
	if err := tmpl.Execute(w, p); err != nil {
		log.Printf("render: %v", err)
	}
}

func handleBg(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "image/webp")
	w.Header().Set("Cache-Control", fmt.Sprintf("public, max-age=%d", 60*60*24*30)) // 30 days
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(bg); err != nil {
		log.Printf("handleBg: %v", err)
	}
}

func main() {
	tinyauthURL := os.Getenv("TINYAUTH_URL")
	demoURL := os.Getenv("DEMO_URL")

	logoutURL := ""

	if tinyauthURL != "" && demoURL != "" {
		encodedDemoURL := url.QueryEscape(strings.TrimSuffix(demoURL, "/") + "/protected")
		logoutURL = strings.TrimSuffix(tinyauthURL, "/") + "/logout?login_for=app&redirect_uri=" + encodedDemoURL
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) { render(w, r, false, logoutURL) })
	mux.HandleFunc("GET /protected", func(w http.ResponseWriter, r *http.Request) { render(w, r, true, logoutURL) })
	mux.HandleFunc("GET /bg.webp", handleBg)

	addr := ":3000"
	log.Printf("tinyauth demo listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
