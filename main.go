package main

import (
	"log"
	"net/http"

	"github.com/gustavokuklinski/hostel-manager/controller/hostel"
	"github.com/gustavokuklinski/hostel-manager/controller/room"
	"github.com/gustavokuklinski/hostel-manager/controller/user"
	"github.com/gustavokuklinski/hostel-manager/lib/loginmanager"
)

func main() {

	Port := "9000"

	log.Println("Listen on: localhost:" + Port)

	// ============================================
	// User Routes [Hostel Admin]
	// ============================================
	// Create/Login user form
	http.HandleFunc("/signin", User.Signin)
	http.HandleFunc("/signup", User.Signup)
	// ============================================
	// Save new user
	http.HandleFunc("/register", User.Register)
	// ============================================
	// Login/Logout
	http.HandleFunc("/login", LoginManager.Login)
	http.HandleFunc("/logout", LoginManager.Logout)
	// ============================================
	// Account setup
	http.HandleFunc("/account/edit", User.EditAcc)
	http.HandleFunc("/account/edit/update", User.UpdateAcc)

	// ============================================
	// Hostel Routes
	// ============================================
	// Template Routes
	http.HandleFunc("/hostel", Hostel.Hostel)
	http.HandleFunc("/hostel/new", Hostel.NewHostel)
	http.HandleFunc("/hostel/edit", Hostel.EditHostel)
	// ===========================================
	// Action Routes
	http.HandleFunc("/hostel/create", Hostel.CreateHostel)
	http.HandleFunc("/hostel/update", Hostel.UpdateHostel)
	http.HandleFunc("/hostel/delete", Hostel.DeleteHostel)
	// ===========================================

	// ===========================================
	// Room Routes
	// ===========================================
	// Template Routes
	http.HandleFunc("/hostel/room", Room.Room)
	http.HandleFunc("/hostel/room/new", Room.NewRoom)
	http.HandleFunc("/hostel/room/edit", Room.EditRoom)
	// ===========================================
	// Action Routes
	http.HandleFunc("/hostel/room/create", Room.CreateRoom)
	http.HandleFunc("/hostel/room/update", Room.UpdateRoom)
	http.HandleFunc("/hostel/room/delete", Room.DeleteRoom)
	// ===========================================

	// ===========================================
	// Reservation Routes
	// ===========================================
	// Template routes
	// http.HandleFunc("/reservation", Reservation.NewReservation)
	// ===========================================
	// Action routes
	// http.HandleFunc("/reservation/create", Reservation.CreateReservation)
	// ===========================================

	// Just to show the Hostels
	// http.HandleFunc("/", Home.Index)

	http.ListenAndServe(":"+Port, nil)
}
