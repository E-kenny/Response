package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/smtp"
	"os"

	"github.com/joho/godotenv"
)

type Response struct {
	Id      int
	Email   string
	Message string
}

func main() {
	godotenv.Load(".env")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if not specified
	}
	server := http.Server{
		Addr: ":" + port,
	}
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))
	http.HandleFunc("/response", response)
	http.HandleFunc("/", home)
	fmt.Println("server is running")
	server.ListenAndServe()
}

func home(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "home.html")
}

func response(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "POST":
		email := r.PostFormValue("email")
		message := r.PostFormValue("message")

		t := template.Must(template.ParseFiles("response.html"))

		emp := Response{}

		emp.Email = email
		emp.Message = message

		send(message, email)

		err := t.Execute(w, emp)
		if err != nil {
			panic(err)
		}

	}
}

//send mail
func send(body, senderEmail string) {
	from := os.Getenv("EMAIL_USERNAME")
	pass := os.Getenv("EMAIL_PASSWORD")
	to := "ekennyobiasogu@gmail.com"

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: PORTFOLIO \n\n" +
		body + senderEmail

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, []string{to}, []byte(msg))

	if err != nil {
		log.Printf("smtp error: %s", err)
		return
	}

	log.Print("sent, check email")
}
