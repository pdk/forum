package srv

import (
	"errors"
	"fmt"
	"html/template"
	"io"
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

// UserError returns the error page with a message for the user.
func (s Server) UserError(w io.Writer, message string) {

	s.WritePage(w, "user-error.html", map[string]interface{}{
		"message": message,
	})
}

// MaybeUserError returns an error page if the condition is true.
func (s Server) MaybeUserError(w io.Writer, condition bool, message string, args ...interface{}) bool {

	if !condition {
		return false
	}

	s.UserError(w, fmt.Sprintf(message, args...))

	return true
}

// bodyAsHTML is a hacky solution to splitting text into paragraphs and maintaining line breaks.
func bodyAsHTML(body string) template.HTML {

	body = strings.ReplaceAll(body, "\r\n", "\n")
	body = strings.ReplaceAll(body, "\r", "\n")
	paragraphs := strings.Split(body, "\n\n")

	sb := strings.Builder{}

	for _, para := range paragraphs {
		para = strings.TrimSpace(para)
		if para == "" {
			continue
		}

		sb.WriteString("<p>\n")
		para = template.HTMLEscapeString(para)
		para = strings.ReplaceAll(para, "\n", "<br>\n")
		sb.WriteString(para)
		sb.WriteString("\n</p>\n")
	}

	return template.HTML(sb.String())
}
