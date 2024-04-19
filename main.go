package main

import (
	"log"
	"net/http"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB

	PORT = "3000"
)

func main() {
	// initialize database
	var err error
	if DB, err = gorm.Open(sqlite.Open("bookstoreapi.db"), &gorm.Config{}); err != nil {
		log.Fatalln("Error connecting to database:", err.Error())
	}

	// run migrations
	if err := DB.AutoMigrate(&Book{}); err != nil {
		log.Fatalln("Error running migrations:", err.Error())
	}

	BookHandler := BookHandler{DB: DB}

	http.HandleFunc("GET /", BookHandler.GetAllBooks)
	http.HandleFunc("GET /create", BookHandler.CreateBook)
	http.HandleFunc("POST /create", BookHandler.StoreBook)
	http.HandleFunc("GET /edit/", BookHandler.EditBook)
	http.HandleFunc("POST /update", BookHandler.UpdateBook)
	http.HandleFunc("POST /delete", BookHandler.DeleteBook)

	log.Printf("Server listening on %q\n", PORT)
	log.Printf("Visit the URL http://127.0.0.1:%s to check if the connection was successful\n", PORT)
	log.Fatalln(http.ListenAndServe(":"+PORT, nil))
}
