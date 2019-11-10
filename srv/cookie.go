package srv

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/pdk/forum/model"
	"github.com/pdk/forum/store"
)

func setCookieValue(w http.ResponseWriter, name, value string) {

	expires := time.Now().AddDate(0, 0, 7)
	cookie := http.Cookie{
		Name:    name,
		Value:   value,
		Expires: expires,
	}

	http.SetCookie(w, &cookie)
}

func getCookieValue(r *http.Request, name string) (string, error) {

	cookie, err := r.Cookie(name)
	if err == http.ErrNoCookie {
		return "", nil
	}

	if err != nil {
		return "", fmt.Errorf("cannot get cookie %s: %w", name, err)
	}

	return cookie.Value, nil
}

func getSignedInUserName(r *http.Request) (string, error) {
	return getCookieValue(r, "name")
}

func setSignedInUserName(w http.ResponseWriter, name string) {
	setCookieValue(w, "name", name)
}

var (
	// ErrNotSignedIn indidcates that there is not a current user.
	ErrNotSignedIn = errors.New("not signed in")
)

// CurrentUser checks the cookie to get current user, and then looks up that
// user in the database, returning the model.User.
func CurrentUser(db *sql.DB, r *http.Request) (model.User, error) {

	userName, err := getSignedInUserName(r)
	if err != nil {
		return model.User{}, fmt.Errorf("cannot get current user: %w", err)
	}

	if userName == "" {
		return model.User{}, fmt.Errorf("%w", ErrNotSignedIn)
	}

	user, err := store.GetUserByName(db, userName)
	if err != nil {
		return model.User{}, fmt.Errorf("cannot get user %s from database: %w", userName, err)
	}

	return user, nil
}
