package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/smtp"
	"os"
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

		t := template.Must(template.ParseFiles("response.html"))

		emp := Response{}

		emp.Email = email
		emp.Message = message

		send(message)

		err := t.Execute(w, emp)
		if err != nil {
			panic(err)
		}

	}
}

//send mail
func send(body string) {
	from := "mathewobiasogu@gmail.com"
	pass := "dempcgxvcdxylohd"
	to := "ekennyobiasogu@gmail.com"

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: Portfolio\n\n" +
		body

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, []string{to}, []byte(msg))

	if err != nil {
		log.Printf("smtp error: %s", err)
		return
	}

	log.Print("sent, check email")
}
