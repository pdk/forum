package srv

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
)

// handleError takes varargs, and assumes the last one is an error var. If not,
// this method will terminate the program (so that we find & fix such coding
// errors).
func handleError(w http.ResponseWriter, message string, args ...interface{}) bool {

	// this will panic if there are 0 args
	lastArg := args[len(args)-1]
	if lastArg == nil {
		return false
	}

	// this will panic if last arg is not actually an error
	err := lastArg.(error)
	if err == nil {
		return false
	}

	message = fmt.Sprintf(message, args...)
	http.Error(w, message, http.StatusInternalServerError)

	return true
}

// errorNotFound checks if the error is a kind of "not found". If so, returns a
// 404 error to client. Returns true to indicate we've already handled the
// client and the page handler should abort processing.
func errorNotFound(w http.ResponseWriter, r *http.Request, err error) bool {

	if errors.Is(err, sql.ErrNoRows) {
		http.NotFound(w, r)
		return true
	}

	return false
}
