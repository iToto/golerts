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
	"strconv"

	"github.com/gorilla/mux"
)

var db *sqlx.DB

type User struct {
	Email    string `db:"email json:"email"`
	Password string `db:"password" json:"password"`
	Token    string `db:"token" json:"token"`
}

type Notification struct {
	Message string `db:"message" json:"message"`
	User    User
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
	router.HandleFunc("/user/{id}/unotifications", ListUserNotifications).Methods("GET")
	router.HandleFunc("/notifications", CreateNotification).Methods("POST")
	log.Fatal(http.ListenAndServe(":"+port, router))
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "OK")
}

func ListUserNotifications(w http.ResponseWriter, r *http.Request) {
	// Get all notifications for user
	vars := mux.Vars(r)
	id, _ := strconv.ParseInt(vars["id"], 10, 0)
	query := `SELECT * FROM notification n
	JOIN user u ON u.email = n.user_email
	WHERE n.user_email = $1`
	notifications := []Notification{}
	db.Select(&notifications, query, id)

	w.Header().Set("Content-Type", "application/json")

	jsonString, _ := json.Marshal(notifications)

	w.Write(jsonString)
}

func GetFoo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseInt(vars["id"], 10, 0)

	var foo string

	err := db.QueryRow("SELECT * FROM foo WHERE id = $1", id).Scan(&foo)

	if err != nil {
		log.Fatal(err)
	}

	// fmt.Print("foo: ", foo)
	fmt.Fprintln(w, "Foo show:", foo)
}

func CreateNotification(w http.ResponseWriter, r *http.Request) {
	// TODO
}
