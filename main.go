package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/E-kenny/resume/print"
	"github.com/go-sql-driver/mysql"
	// _ "github.com/go-sql-driver/mysql"
)

type Response struct {
	Id      int
	Email   string
	Message string
}

func main() {
	//Task 2, this functon prints my name
	print.Name()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if not specified
	}
	server := http.Server{
		Addr: ":" + port,
	}
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))
	http.HandleFunc("/resume", resume)
	http.HandleFunc("/", home)
	fmt.Println("server is running")
	server.ListenAndServe()
}

func home(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "home.html")
}

func resume(w http.ResponseWriter, r *http.Request) {

	switch r.Method {

	case "GET":
		http.ServeFile(w, r, "resume.html")
	case "POST":
		email := r.PostFormValue("email")
		message := r.PostFormValue("message")

		db := dbConn()
		_, err := db.Exec("INSERT INTO responses (email, message) VALUES ( ?, ?)", email, message)

		if err != nil {
			fmt.Printf("Addresponse: %v", err)
		}

		t := template.Must(template.ParseFiles("response.html"))

		selDB, err := db.Query("SELECT * FROM responses WHERE email=?", email)
		if err != nil {
			panic(err.Error())
		}
		emp := Response{}

		for selDB.Next() {
			var id int
			var email, message string
			err = selDB.Scan(&id, &email, &message)
			if err != nil {
				panic(err.Error())
			}
			emp.Id = id
			emp.Email = email
			emp.Message = message
		}

		err = t.Execute(w, emp)
		if err != nil {
			panic(err)
		}
	}
}

func dbConn() (db *sql.DB) {

	cfg := mysql.Config{
		User:   os.Getenv("DB_USERNAME"),
		Passwd: os.Getenv("DB_PASSWORD"),
		Net:    "tcp",
		Addr:   os.Getenv("DB_HOST"),
		DBName: os.Getenv("DB_DATABASE"),
	}
	dbDriver := os.Getenv("DB_CONNECTION")
	cfg.AllowNativePasswords = true

	db, err := sql.Open(dbDriver, cfg.FormatDSN())
	if err != nil {
		panic(err.Error())
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")
	return db
}
