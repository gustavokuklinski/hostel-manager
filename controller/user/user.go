package User

import (
	"log"
	"net/http"
	"strconv"
	"text/template"

	"github.com/gustavokuklinski/hoxtel.me/lib/crypty"
	"github.com/gustavokuklinski/hoxtel.me/lib/loginmanager"
	"github.com/gustavokuklinski/hoxtel.me/lib/sysdb"
)

type sUser struct {
	Cid, Id                        int
	Name, Surname, Email, Password string
}

var tmpl = template.Must(template.ParseGlob("view/user/*"))

// Signin Form
func Signin(w http.ResponseWriter, r *http.Request) {
	log.Println("Sign In")
	tmpl.ExecuteTemplate(w, "Signin", nil)
}

// Signup Form
func Signup(w http.ResponseWriter, r *http.Request) {
	log.Println("Sign Up")
	tmpl.ExecuteTemplate(w, "Signup", nil)
}

// Signin action to register new user
func Register(w http.ResponseWriter, r *http.Request) {
	db := sysdb.DbConn()

	if r.Method == "POST" {
		name := r.FormValue("name")
		surname := r.FormValue("surname")
		email := r.FormValue("email")
		password := r.FormValue("password")

		registerUser, err := db.Prepare("INSERT INTO users(name, surname, email, password) VALUES(?,?,?,?)")
		if err != nil {
			panic(err.Error())
		}

		registerUser.Exec(name, surname, email, password)

		log.Println("Register new user: e-mail: " + email + " | password: " + password)
	}
	defer db.Close()
	http.Redirect(w, r, "/signin", 301)
}

// Edit user account
func EditAcc(w http.ResponseWriter, r *http.Request) {
	db := sysdb.DbConn()

	LoginManager.CheckActive(w, r)

	getId := r.URL.Query().Get("uid")
	StringInt, _ := strconv.Atoi(getId)
	hashId := Crypty.Decrypto(StringInt)
	hashString := strconv.Itoa(hashId)

	selectUser, err := db.Query("SELECT id, name, surname, email, password FROM users WHERE id=" + hashString)
	if err != nil {
		panic(err.Error())
	}

	user := sUser{}

	for selectUser.Next() {
		var id int
		var name, surname, email, password string

		err = selectUser.Scan(&id, &name, &surname, &email, &password)
		if err != nil {
			panic(err.Error())
		}

		user.Id = id
		user.Name = name
		user.Surname = surname
		user.Email = email
		user.Password = password
		user.Cid = Crypty.Crypto(user.Id)
	}

	tmpl.ExecuteTemplate(w, "Edit", user)

	log.Println("Editing user account:" + user.Name + " " + user.Surname + " | Email: " + user.Email)

	defer db.Close()
}

// Update user account
func UpdateAcc(w http.ResponseWriter, r *http.Request) {
	db := sysdb.DbConn()

	LoginManager.CheckActive(w, r)

	if r.Method == "POST" {
		uid := r.FormValue("uid")
		StringId, _ := strconv.Atoi(uid)
		hashId := Crypty.Decrypto(StringId)
		hashString := strconv.Itoa(hashId)

		name := r.FormValue("name")
		surname := r.FormValue("surname")
		email := r.FormValue("email")
		password := r.FormValue("password")

		updateUser, err := db.Prepare("UPDATE users SET name=?, surname=?, email=?, password=? WHERE id=?")
		if err != nil {
			panic(err.Error())
		}

		updateUser.Exec(name, surname, email, password, hashString)

		log.Println("User updated...Loggin Out")

		LoginManager.Logout(w, r)
	}

	db.Close()

	http.Redirect(w, r, "/hostel", 301)
}
