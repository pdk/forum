package srv

import (
	"database/sql"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
)

// Server handles incoming HTTP requests.
type Server struct {
	DB        *sql.DB
	AssetsDir string
	Template  *template.Template
}

// NewServer construct and return a new Server.
func NewServer(db *sql.DB, assetsDir string) (Server, error) {

	templateGlob := assetsDir + "/templates/*.html"
	log.Printf("reading & parsing templates in %s", templateGlob)

	tmpl, err := template.ParseGlob(templateGlob)
	if err != nil {
		return Server{}, fmt.Errorf("failed to compile templates from %s: %w", templateGlob, err)
	}

	for _, t := range tmpl.Templates() {
		log.Printf("template ready: %s", t.Name())
	}

	return Server{
		DB:        db,
		AssetsDir: assetsDir,
		Template:  tmpl,
	}, nil
}

// ListenAndServe sets up routes and kicks off HTTP listener.
func (s Server) ListenAndServe(listenAddress string) {

	static := http.FileServer(http.Dir(s.AssetsDir + "/static/"))
	log.Printf("routing /css/, /js/, /img/ => %v", static)
	http.Handle("/css/", static)
	http.Handle("/js/", static)
	http.Handle("/img/", static)

	routes := map[string]http.HandlerFunc{
		// just using map to make formatting easier to read
		"/":           s.HomePage,
		"/sign-in":    s.SignIn,
		"/topics":     s.OnlySignedIn(s.TopicsPage),
		"/add-topic":  s.OnlySignedIn(s.AddTopic),
		"/topics/":    s.OnlySignedIn(s.OneTopicPage),
		"/add-thread": s.OnlySignedIn(s.AddThread),
		"/threads/":   s.OnlySignedIn(s.OneThreadPage),
		"/add-post":   s.OnlySignedIn(s.AddPost),
	}

	for path, handler := range routes {
		log.Printf("routing %s => %v", path, handler)
		http.HandleFunc(path, handler)
	}

	log.Printf("listening at %s", listenAddress)
	log.Fatalf("server failed: %w",
		http.ListenAndServe(listenAddress, nil))
}

// WritePage executes a named template with the given data. This is meant to be
// called by a page handler, and as the last thing done by page handlers,
// there's nowhere to send an error, so we just log any errors here.
func (s Server) WritePage(w io.Writer, name string, data interface{}) {

	err := s.Template.ExecuteTemplate(w, name, data)

	if err != nil {
		log.Printf("error executing template %s: %w", name, err)
	}
}
