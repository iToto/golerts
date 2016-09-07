package main

import (
	// "database/sql"
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

var db *sqlx.DB

type User struct {
	Id       string `db:"id json:"id"`
	Email    string `db:"email json:"email"`
	Password string `db:"password" json:"password"`
}

type Notification struct {
	Id      string `db:"id json:"id"`
	Message string `db:"message" json:"message"`
	User    User
}

type Token struct {
	Id     string `db:"id json:"id"`
	Token  string `db:"token json:"token"`
	UserId string `db:"user_id json:"user_id"`
	Status bool   `db:"status json:"status"`
}

func DBConnection() *sqlx.DB {
	dbURL := os.Getenv("DATABASE_URL")

	if dbURL == "" {
		log.Fatal("$DATABASE_URL environment variable must be set")
	} else {
		fmt.Println("Connected to database: " + dbURL)
	}

	db, err := sqlx.Connect("postgres", dbURL)

	if err != nil {
		log.Fatal(err)
	}

	return db
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT environment variable must be set")
	} else {
		fmt.Println("Running on port: " + port)
	}

	db = DBConnection()
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", Index)
	router.HandleFunc("/user/{id:[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}}/notifications", ListUserNotifications).Methods("GET")
	router.HandleFunc("/notifications", CreateNotification).Methods("POST")
	log.Fatal(http.ListenAndServe(":"+port, router))
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "OK")
}

func ListUserNotifications(w http.ResponseWriter, r *http.Request) {
	// Get all notifications for user
	vars := mux.Vars(r)
	id, _ := vars["id"]
	query := `SELECT n.id, n.message, u.id as user_id, u.email as user_email, u.password as user_password FROM notification n
	JOIN "user" u ON u.id = n.user_id
	WHERE n.user_id = '$1'`

	fmt.Printf("User ID %v", id)

	notifications := []Notification{}

	db.Select(&notifications, query, id)
	fmt.Printf("Notifications: %v", notifications)
	w.Header().Set("Content-Type", "application/json")

	jsonString, _ := json.Marshal(notifications)

	w.Write(jsonString)
}

func CreateNotification(w http.ResponseWriter, r *http.Request) {
	// TODO
}
