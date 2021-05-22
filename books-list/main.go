package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Book struct {
	Id     int    `json:id`
	Title  string `json:title`
	Author string `json:author`
	Year   string `json:year`
}

var books []Book

func main() {

	router := mux.NewRouter()
	fmt.Println("server running...")

	books = append(books, Book{Id: 1, Title: "Golang pointers", Author: "Mr Motke", Year: "2000"},
		Book{Id: 2, Title: "Golang Api", Author: "Mr Gholap", Year: "1999"},
		Book{Id: 3, Title: "Golang Concurrency", Author: "Mr Vats", Year: "1995"},
		Book{Id: 4, Title: "Golang Deploy", Author: "Mr Gupta", Year: "1997"},
		Book{Id: 5, Title: "Golang Good", Author: "Mr Saini", Year: "1998"})

	router.HandleFunc("/books", getBooks).Methods("GET")
	router.HandleFunc("/books/{id}", getBook).Methods("GET")
	router.HandleFunc("/books/add", addBook).Methods("POST")
	router.HandleFunc("/books", updateBook).Methods("PUT")
	router.HandleFunc("/books/{id}", removeBook).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8001", router))
	fmt.Println("running on port 8001...")

}

func getBooks(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(books)
}

func getBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	log.Println(params)

	i, _ := strconv.Atoi(params["id"])
	for _, book := range books {
		if book.Id == i {
			json.NewEncoder(w).Encode(&book)
		}
	}
}

func addBook(w http.ResponseWriter, r *http.Request) {
	log.Println("Add Book is called")
	var newbook Book
	_ = json.NewDecoder(r.Body).Decode(&newbook)

	books = append(books, newbook)

	json.NewEncoder(w).Encode(books)

	//json.NewEncoder(w).Encode("New Book Added to library")
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	log.Println("Update Book is called")

	// var upbook Book

	// json.NewDecoder(r.Body).Decode(&books)

	// for i, item := range books {
	// 	if item.Id == Book.id {
	// 		books[i] = upbooks
	// 	}
	// }
	var bookUpdate Book
	json.NewDecoder(r.Body).Decode(&bookUpdate)

	var result []Book

	for _, item := range books {
		if item.Id == bookUpdate.Id {

			if bookUpdate.Title != "" {
				item.Title = bookUpdate.Title
			}
			if bookUpdate.Author != "" {
				item.Author = bookUpdate.Author
			}
			if bookUpdate.Year != "" {
				item.Year = bookUpdate.Year
			}

		}
		result = append(result, item)
	}
	json.NewEncoder(w).Encode(result)
}

func removeBook(w http.ResponseWriter, r *http.Request) {
	log.Println("Remove Book is called")

	paams := mux.Vars(r)

	id, _ := strconv.Atoi(paams["id"])

	for i, item := range books {
		if item.Id == id {
			books = append(books[:i], books[i+1:]...)
		}
	}

	json.NewEncoder(w).Encode(books)
}
