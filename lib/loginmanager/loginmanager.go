package LoginManager

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/gustavokuklinski/hoxtel.me/lib/sysdb"
)

var store = sessions.NewCookieStore([]byte("9208439023890284028490"))

func Login(w http.ResponseWriter, r *http.Request) {
	db := sysdb.DbConn()

	if r.Method == "POST" {
		formemail := r.FormValue("email")
		formpassword := r.FormValue("password")

		var id int
		var email, password string

		// Execute query looking for form parameters(Email and Password)
		err := db.QueryRow("SELECT id, email, password FROM users WHERE email=? AND password=?", formemail, formpassword).Scan(&id, &email, &password)

		// This switch will check for valid rows
		switch {
		// Case invalid, log error and redirect to login
		case err == sql.ErrNoRows:
			log.Println("Error loggin in")
			http.Redirect(w, r, "/signin", 301)

			// Case wrong query display panic on console
		case err != nil:
			panic(err.Error())

			// Default as Success
		default:
			// Create session using the e-mail
			session, _ := store.Get(r, "uid-hostel-session")
			session.Values["uid-user"] = email
			session.Options = &sessions.Options{MaxAge: 3600 * 1, HttpOnly: true} // Set the session to 1 hour and only HTTP
			// Save the session
			session.Save(r, w)

			// Log on console and redirect to hostel page
			log.Println("Succefull login")
			http.Redirect(w, r, "/hostel", 301)
		}
	}

	defer db.Close()

}

func Logout(w http.ResponseWriter, r *http.Request) {
	// Get the session ID
	session, _ := store.Get(r, "uid-hostel-session")

	// Turn the session age to negative, killing it
	session.Options.MaxAge = -1

	// Save the new session, loggin out
	session.Save(r, w)

	// Log screen and redirect to login form
	log.Println("Logged out")
	http.Redirect(w, r, "/signin", 301)

}

// Verify if the session is alive
func CheckActive(w http.ResponseWriter, r *http.Request) {
	loginSess, _ := store.Get(r, "uid-hostel-session")

	if _, ok := loginSess.Values["uid-user"]; !ok {
		log.Println("User not logged in")
		http.Redirect(w, r, "/signin", 301)
	}
}

func GetActive(r *http.Request) string {
	loginSess, _ := store.Get(r, "uid-hostel-session")

	return loginSess.Values["uid-user"].(string)

}

func GetActiveBool(r *http.Request) interface{} {
	loginSess, _ := store.Get(r, "uid-hostel-session")

	return loginSess.Values["uid-user"]

}
