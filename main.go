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
	apns "github.com/sideshow/apns2"
	"github.com/sideshow/apns2/certificate"
)

var db *sqlx.DB
var apnsClient *apns.Client

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

const banner = `
(_______|_______|_)     (_______|_____ (_______) _____)
 _   ___ _     _ _       _____   _____) )  _  ( (____
| | (_  | |   | | |     |  ___) |  __  /  | |  \____ \
| |___) | |___| | |_____| |_____| |  \ \  | |  _____) )
 \_____/ \_____/|_______)_______)_|   |_| |_| (______/
`

func DBConnection() *sqlx.DB {
	dbURL := os.Getenv("DATABASE_URL")

	if dbURL == "" {
		log.Fatal("$DATABASE_URL environment variable must be set")
	} else {
		log.Println("Connected to database: " + dbURL)
	}

	db, err := sqlx.Connect("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func APNSConnection() *apns.Client {
	// TODO: This should be done through configuration
	cert, err := certificate.FromP12File("devPush.p12", "aaaaaaaaAA$%")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to APNS")
	return apns.NewClient(cert).Production()
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT environment variable must be set")
	} else {
		fmt.Println("Running on port: " + port)
	}

	db = DBConnection()
	apnsClient = APNSConnection()
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", Index)
	router.HandleFunc("/login", UserLogin).Methods("POST")
	router.HandleFunc("/tokens", CreateToken).Methods("POST")
	router.HandleFunc("/notifications", CreateNotification).Methods("POST")
	router.HandleFunc("/users/{id:[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}}/notifications", ListUserNotifications).Methods("GET")

	log.Fatal(http.ListenAndServe(":"+port, router))
}

/// UserLogin will register a user with the service so that they can receive notifications
func UserLogin(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var user User

	err := decoder.Decode(&user)
	if err != nil {
		log.Fatal(err)
	}

	id, err := uuid.NewV4()
	if err != nil {
		log.Fatal(err)
	}

	user.Id = id.String()
	query := "INSERT INTO \"user\" (id, email, password) VALUES (:id, :email, :password)"
	_, err = db.NamedExec(query, &user)
	if err != nil {
		log.Fatal(err)
	}

	jsonString, _ := json.Marshal(user)

	w.Write(jsonString)
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
	decoder := json.NewDecoder(r.Body)

	var notification Notification

	err := decoder.Decode(&notification)
	if err != nil {
		log.Fatal(err)
	}

	// TODO: Get user's token
	token := getActiveTokenForUserId(notification.UserId)

	// Send notification to APNS server
	apnsNotification := &apns.Notification{}
	apnsNotification.DeviceToken = token.Token
	apnsNotification.Topic = "com.iToto.golerts"
	apnsNotification.Payload = []byte(`{"aps":{"alert":"Hello World!"}}`)

	res, err := apnsClient.Push(apnsNotification)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("APNs ID:", res.ApnsID)

	u, err := uuid.NewV4()
	if err != nil {
		log.Fatal(err)
	}

	notification.Id = u.String()
	query := "INSERT INTO notification (id, message, user_id) VALUES (:id, :message, :user_id)"
	_, err = db.NamedExec(query, &notification)
	if err != nil {
		log.Fatal(err)
	}

	jsonString, _ := json.Marshal(notification)

	// TODO: Make call to Apple Servers to broadcast notification to user's device

	w.Write(jsonString)

}

func CreateToken(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var token Token

	err := decoder.Decode(&token)
	if err != nil {
		log.Fatal(err)
	}

	id, err := uuid.NewV4()
	if err != nil {
		log.Fatal(err)
	}

	token.Id = id.String()
	query := "INSERT INTO token (id, token, user_id, status) VALUES (:id, :token, :user_id, :status)"
	_, err = db.NamedExec(query, &token)
	if err != nil {
		log.Fatal(err)
	}

	jsonString, _ := json.Marshal(token)

	w.Write(jsonString)
}

func getUserById(userId string) *User {
	var user User
	query := "SELECT * FROM \"user\" WHERE id = $1"

	err := db.Get(&user, query, userId)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("User with id: %v", user)

	return &user
}

func getActiveTokenForUserId(userId string) *Token {
	var token Token
	query := "SELECT * FROM token WHERE user_id = $1 AND status = TRUE"

	err := db.Get(&token, query, userId)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Token for user", token)

	return &token
}
