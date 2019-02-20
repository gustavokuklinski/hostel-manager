package Room

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
type sRoom struct {
	Cid, Id, Hostel_id, Capacity, Locker, Bathroom, Couple_bad, Additional_bad int
	Name                                                                       string
}
type sUser struct {
	Cid, Id int
}

var tmpl = template.Must(template.ParseGlob("view/room/*"))

func Room(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/hostel", 301)
}

func NewRoom(w http.ResponseWriter, r *http.Request) {
	db := sysdb.DbConn()

	LoginManager.CheckActive(w, r)

	if LoginManager.GetActiveBool(r) != nil {

		// If ok, redirect to hostel page
		log.Println("Create new Room")

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

		getHostelId := r.URL.Query().Get("huid")
		StringInt, _ := strconv.Atoi(getHostelId)
		hashId := Crypty.Decrypto(StringInt)
		hashString := strconv.Itoa(hashId)

		selectHostel, err := db.Query("SELECT id FROM hostel WHERE id=" + hashString)
		if err != nil {
			panic(err.Error())
		}

		Hostel := sHostel{}

		for selectHostel.Next() {
			var id int
			err = selectHostel.Scan(&id)
			if err != nil {
				panic(err.Error())
			}
			Hostel.Id = id
			Hostel.Cid = Crypty.Crypto(Hostel.Id)
		}

		if err := tmpl.ExecuteTemplate(w, "New",
			struct{ User, Hostel interface{} }{User, Hostel}); err != nil {
			//
		}
	}
	db.Close()
}

func EditRoom(w http.ResponseWriter, r *http.Request) {
	db := sysdb.DbConn()

	LoginManager.CheckActive(w, r)
	if LoginManager.GetActiveBool(r) != nil {

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

		getRoomId := r.URL.Query().Get("ruid")
		StringInt, _ := strconv.Atoi(getRoomId)
		hashId := Crypty.Decrypto(StringInt)
		hashString := strconv.Itoa(hashId)

		selectRoom, err := db.Query("SELECT id, name, capacity, locker, bathroom, couple_bad, additional_bad FROM room WHERE id=" + hashString)
		if err != nil {
			panic(err.Error())
		}

		Room := sRoom{}

		for selectRoom.Next() {
			var id, capacity, locker, bathroom, couple_bad, additional_bad int
			var name string
			err = selectRoom.Scan(&id, &name, &capacity, &locker, &bathroom, &couple_bad, &additional_bad)
			if err != nil {
				panic(err.Error())
			}

			Room.Id = id
			Room.Name = name
			Room.Capacity = capacity
			Room.Locker = locker
			Room.Bathroom = bathroom
			Room.Couple_bad = couple_bad
			Room.Additional_bad = additional_bad

			Room.Cid = Crypty.Crypto(Room.Id)
		}

		if err := tmpl.ExecuteTemplate(w, "Edit",
			struct{ Room, User interface{} }{Room, User}); err != nil {
			//
		}
	}
	db.Close()

}

func CreateRoom(w http.ResponseWriter, r *http.Request) {
	db := sysdb.DbConn()

	LoginManager.CheckActive(w, r)

	if LoginManager.GetActiveBool(r) != nil {
		if r.Method == "POST" {
			huid := r.FormValue("huid")
			StringInt, _ := strconv.Atoi(huid)
			hashId := Crypty.Decrypto(StringInt)
			hashString := strconv.Itoa(hashId)

			name := r.FormValue("name")
			capacity := r.FormValue("capacity")
			locker := r.FormValue("locker")
			bathroom := r.FormValue("bathroom")
			coupleBad := r.FormValue("couple_bad")
			additionalBad := r.FormValue("additional_bad")

			registerRoom, err := db.Prepare("INSERT INTO room(name, capacity, locker, bathroom, couple_bad, additional_bad, hostel_id) VALUES(?,?,?,?,?,?,?)")
			if err != nil {
				panic(err.Error())
			}

			registerRoom.Exec(name, capacity, locker, bathroom, coupleBad, additionalBad, hashString)
			log.Println("New Room on database: " + name)
		}

		db.Close()
	}
	http.Redirect(w, r, "/hostel", 301)
}

func UpdateRoom(w http.ResponseWriter, r *http.Request) {
	db := sysdb.DbConn()

	LoginManager.CheckActive(w, r)
	if LoginManager.GetActiveBool(r) != nil {
		if r.Method == "POST" {
			ruid := r.FormValue("ruid")
			StringInt, _ := strconv.Atoi(ruid)
			hashId := Crypty.Decrypto(StringInt)
			hashString := strconv.Itoa(hashId)

			name := r.FormValue("name")
			capacity := r.FormValue("capacity")
			locker := r.FormValue("locker")
			bathroom := r.FormValue("bathroom")
			coupleBad := r.FormValue("couple_bad")
			additionalBad := r.FormValue("additional_bad")

			updateRoom, err := db.Prepare("UPDATE room SET name=?, capacity=?, locker=?, bathroom=?, couple_bad=?, additional_bad=? WHERE id=?")
			if err != nil {
				panic(err.Error())
			}

			updateRoom.Exec(name, capacity, locker, bathroom, coupleBad, additionalBad, hashString)
			log.Println("Room updated")
		}
	}
	db.Close()
	http.Redirect(w, r, "/hostel", 301)
}

func DeleteRoom(w http.ResponseWriter, r *http.Request) {
	db := sysdb.DbConn()

	LoginManager.CheckActive(w, r)
	if LoginManager.GetActiveBool(r) != nil {
		getRoomId := r.URL.Query().Get("ruid")
		StringInt, _ := strconv.Atoi(getRoomId)
		hashId := Crypty.Decrypto(StringInt)
		hashString := strconv.Itoa(hashId)

		deleteRoom, err := db.Prepare("DELETE FROM room WHERE id=?")
		if err != nil {
			panic(err.Error())
		}
		deleteRoom.Exec(hashString)

		log.Println("Room deleted")

		db.Close()
	}
	http.Redirect(w, r, "/hostel", 301)
}
