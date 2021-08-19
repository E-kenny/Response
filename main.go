package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type Response struct {
	Id      int
	Email   string
	Message string
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if not specified
	}
	server := http.Server{
		Addr: ":" + port,
	}
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))
	http.HandleFunc("/", resume)
	fmt.Println("server is running")
	server.ListenAndServe()
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
	// os.Setenv("USERNAME", "root")
	// os.Setenv("PASSWORD", "E_kenny246810")
	dbDriver := "mysql"
	dbUser := os.Getenv("USERNAME") //"root"
	dbPass := os.Getenv("PASSWORD") //"E_kenny246810"
	dbName := "go-mysql-crud"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}
