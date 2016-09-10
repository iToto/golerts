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
	"github.com/nu7hatch/gouuid"
)

var db *sqlx.DB

type User struct {
	Id       string `json:"id" db:"id"`
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"password"`
}

type Notification struct {
	Id      string `json:"id" db:"id"`
	Message string `json:"message" db:"message"`
	UserId  string `json:"user_id" db:"user_id"`
}

type Token struct {
	Id     string `json:"id" db:"id"`
	Token  string `json:"token" db:"token"`
	UserId string `json:"user_id" db:"user_id"`
	Status bool   `json:"status" db:"status"`
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
	query := `SELECT id, message, user_id FROM notification WHERE user_id = $1`

	fmt.Printf("User ID %v", id)

	var notifications []*Notification

	db.Select(&notifications, query, id)
	fmt.Printf("Notifications: %v", notifications)
	w.Header().Set("Content-Type", "application/json")

	jsonString, _ := json.Marshal(notifications)

	w.Write(jsonString)
}

func CreateNotification(w http.ResponseWriter, r *http.Request) {
	// TODO
	decoder := json.NewDecoder(r.Body)

	var notification Notification

	err := decoder.Decode(&notification)
	if err != nil {
		log.Fatal(err)
	}

	u, err := uuid.NewV4()
	if err != nil {
		log.Fatal(err)
	}

	notification.Id = u.String()
	query := "INSERT INTO notification (id, message, user_id) VALUES (:id, :message, :user_id)"
	result, err := db.NamedExec(query, &notification)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Result: %v", result)
	jsonString, _ := json.Marshal(notification)

	w.Write(jsonString)

}
