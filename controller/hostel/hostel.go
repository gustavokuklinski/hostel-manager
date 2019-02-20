package Hostel

import (
	"log"
	"net/http"
	"strconv"
	"text/template"

	"github.com/gustavokuklinski/hoxtel.me/lib/crypty"
	"github.com/gustavokuklinski/hoxtel.me/lib/loginmanager"
	"github.com/gustavokuklinski/hoxtel.me/lib/sysdb"
)

type sHostel struct {
	Cid, Id, User_id            int
	Name, State, Address, Phone string
}
type sUser struct {
	Cid, Id int
}
type sRoom struct {
	Cid, Id, Hostel_id, Capacity, Locker, Bathroom, Couple_bad, Additional_bad int
	Name                                                                       string
}

var tmpl = template.Must(template.ParseGlob("view/hostel/*"))

func Hostel(w http.ResponseWriter, r *http.Request) {
	db := sysdb.DbConn()

	LoginManager.CheckActive(w, r)

	if LoginManager.GetActiveBool(r) != nil {

		// If ok, redirect to hostel page
		log.Println("Listing all hostel - User logged in: " + LoginManager.GetActive(r))

		selectUser, err := db.Query("SELECT id FROM users WHERE email='" + LoginManager.GetActive(r) + "'")
		if err != nil {
			panic(err.Error())
		}

		User := sUser{}

		for selectUser.Next() {
			var id int
			err = selectUser.Scan(&id)
			if err != nil {
				panic(err.Error())
			}
			User.Id = id
			User.Cid = Crypty.Crypto(User.Id)
		}

		strId := User.Id

		selectHostel, err := db.Query("SELECT id, name FROM hostel WHERE user_id=" + strconv.Itoa(strId))
		if err != nil {
			panic(err.Error())
		}
		Hostel := sHostel{}

		for selectHostel.Next() {
			var id int
			var name string
			err = selectHostel.Scan(&id, &name)
			if err != nil {
				panic(err.Error())
			}
			Hostel.Id = id
			Hostel.Name = name
			Hostel.Cid = Crypty.Crypto(Hostel.Id)
		}

		Room := sRoom{}
		SliceRoom := []sRoom{}

		strHuid := Hostel.Id

		selectRoom, err := db.Query("SELECT id, name FROM room WHERE hostel_id=" + strconv.Itoa(strHuid))
		if err != nil {
			panic(err.Error())
		}

		for selectRoom.Next() {
			var id int
			var name string
			err = selectRoom.Scan(&id, &name)
			if err != nil {
				panic(err.Error())
			}
			Room.Id = id
			Room.Name = name
			Room.Cid = Crypty.Crypto(Room.Id)

			SliceRoom = append(SliceRoom, Room)
		}

		if err := tmpl.ExecuteTemplate(w, "Hostel",
			struct{ User, Hostel, SliceRoom interface{} }{User, Hostel, SliceRoom}); err != nil {
			//
		}
	}

	defer db.Close()
}

func NewHostel(w http.ResponseWriter, r *http.Request) {
	db := sysdb.DbConn()

	LoginManager.CheckActive(w, r)

	if LoginManager.GetActiveBool(r) != nil {

		log.Println("Create new Hostel")

		selectUser, err := db.Query("SELECT id FROM users WHERE email='" + LoginManager.GetActive(r) + "'")
		if err != nil {
			panic(err.Error())
		}

		User := sUser{}

		for selectUser.Next() {
			var id int
			err = selectUser.Scan(&id)
			if err != nil {
				panic(err.Error())
			}
			User.Id = id
			User.Cid = Crypty.Crypto(User.Id)

		}

		if err := tmpl.ExecuteTemplate(w, "New",
			struct{ User interface{} }{User}); err != nil {
		}

		db.Close()
	}
}

func EditHostel(w http.ResponseWriter, r *http.Request) {
	db := sysdb.DbConn()

	LoginManager.CheckActive(w, r)

	if LoginManager.GetActiveBool(r) != nil {

		getHostelId := r.URL.Query().Get("huid")
		StringInt, _ := strconv.Atoi(getHostelId)
		hashId := Crypty.Decrypto(StringInt)
		hashString := strconv.Itoa(hashId)
		selectUser, err := db.Query("SELECT id FROM users WHERE email='" + LoginManager.GetActive(r) + "'")

		if err != nil {
			panic(err.Error())
		}

		User := sUser{}

		for selectUser.Next() {
			var id int
			err = selectUser.Scan(&id)
			if err != nil {
				panic(err.Error())
			}
			User.Id = id
			User.Cid = Crypty.Crypto(User.Id)
		}

		selectHostel, err := db.Query("SELECT id, name, state, address, phone FROM hostel WHERE id=" + hashString)
		if err != nil {
			panic(err.Error())
		}

		Hostel := sHostel{}

		for selectHostel.Next() {
			var id int
			var name, state, address, phone string
			err = selectHostel.Scan(&id, &name, &state, &address, &phone)
			if err != nil {
				panic(err.Error())
			}

			Hostel.Id = id
			Hostel.Name = name
			Hostel.State = state
			Hostel.Address = address
			Hostel.Phone = phone
			Hostel.Cid = Crypty.Crypto(Hostel.Id)

		}

		if err := tmpl.ExecuteTemplate(w, "Edit",
			struct{ Hostel, User interface{} }{Hostel, User}); err != nil {
		}

		db.Close()
	}
}

func CreateHostel(w http.ResponseWriter, r *http.Request) {
	db := sysdb.DbConn()

	LoginManager.CheckActive(w, r)

	if LoginManager.GetActiveBool(r) != nil {
		if r.Method == "POST" {
			name := r.FormValue("name")
			state := r.FormValue("state")
			address := r.FormValue("address")
			phone := r.FormValue("phone")
			user_id := r.FormValue("uid")
			StringInt, _ := strconv.Atoi(user_id)
			hashId := Crypty.Decrypto(StringInt)
			hashString := strconv.Itoa(hashId)

			registerHostel, err := db.Prepare("INSERT INTO hostel(name, state, address, phone, user_id) VALUES(?,?,?,?,?)")
			if err != nil {
				panic(err.Error())
			}

			registerHostel.Exec(name, state, address, phone, hashString)
			log.Println("New hostel on database: " + name)
		}

		db.Close()
	}
	http.Redirect(w, r, "/hostel", 301)
}

func UpdateHostel(w http.ResponseWriter, r *http.Request) {
	db := sysdb.DbConn()

	LoginManager.CheckActive(w, r)

	if LoginManager.GetActiveBool(r) != nil {
		if r.Method == "POST" {
			name := r.FormValue("name")
			state := r.FormValue("state")
			address := r.FormValue("address")
			phone := r.FormValue("phone")
			getHostelId := r.FormValue("huid")

			StringInt, _ := strconv.Atoi(getHostelId)
			hashId := Crypty.Decrypto(StringInt)
			hashString := strconv.Itoa(hashId)

			updateHostel, err := db.Prepare("UPDATE hostel SET name=?, state=?, address=?, phone=? WHERE id=?")
			if err != nil {
				panic(err.Error())
			}
			updateHostel.Exec(name, state, address, phone, hashString)
			log.Println("Hostel updated")
		}

		db.Close()

		http.Redirect(w, r, "/hostel", 301)
	}
}

func DeleteHostel(w http.ResponseWriter, r *http.Request) {
	db := sysdb.DbConn()

	LoginManager.CheckActive(w, r)

	if LoginManager.GetActiveBool(r) != nil {
		getHostelId := r.URL.Query().Get("huid")
		StringInt, _ := strconv.Atoi(getHostelId)
		hashId := Crypty.Decrypto(StringInt)
		hashString := strconv.Itoa(hashId)

		deleteHostel, err := db.Prepare("DELETE FROM hostel WHERE id=?")
		if err != nil {
			panic(err.Error())
		}
		deleteHostel.Exec(hashString)

		log.Println("Hostel deleted")

		db.Close()
	}

	http.Redirect(w, r, "/hostel", 301)
}
