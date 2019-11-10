package srv

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// OnlySignedIn will redirect to front page if the user is not signed in.
func (s Server) OnlySignedIn(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		userName, err := getSignedInUserName(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		if userName == "" {
			// not signed in
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}

		handler(w, r)
	}
}

func getPathID(url *url.URL) (int64, error) {

	uriParts := strings.Split(url.Path, "/")
	if len(uriParts) < 3 {
		return 0, errors.New("no ID present")
	}

	id, err := strconv.ParseInt(uriParts[2], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("cannot parse ID %s: %w", uriParts[2], err)
	}

	return id, nil
}

// ErrorPage returns the error page with a message for the user.
func (s Server) ErrorPage(w http.ResponseWriter, message string) {
	err := s.Template.ExecuteTemplate(w, "user-error.html", map[string]interface{}{
		"message": message,
	})
	if err != nil {
		log.Printf("error executing template user-error.html: %w", err)
	}
}
