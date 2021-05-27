package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/lib/pq"
	"github.com/subosito/gotenv"
)

type Book struct {
	Id     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
	Year   string `json:"year"`
}

var books []Book

var db *sql.DB

func init() {

	gotenv.Load()
}

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {

	pgUrl, err := pq.ParseURL(os.Getenv("ElephantSQL_URL"))
	logFatal(err)

	db, err = sql.Open("postgres", pgUrl)
	logFatal(err)

	err = db.Ping()
	logFatal(err)

	log.Println(pgUrl)

	router := mux.NewRouter()
	fmt.Println("server running on port 9000...")

	router.HandleFunc("/books", getBooks).Methods("GET")
	router.HandleFunc("/books/{id}", getBook).Methods("GET")
	router.HandleFunc("/books/add", addBook).Methods("POST")
	router.HandleFunc("/books", updateBook).Methods("PUT")
	router.HandleFunc("/books/{id}", removeBook).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":9000", router))

}

func getBooks(w http.ResponseWriter, r *http.Request) {
	var book Book
	books = []Book{}

	rows, err := db.Query("Select * from books")
	logFatal(err)

	for rows.Next() {
		err := rows.Scan(&book.Id, &book.Title, &book.Author, &book.Year)
		logFatal(err)

		books = append(books, book)
	}

	json.NewEncoder(w).Encode(books)
}

func getBook(w http.ResponseWriter, r *http.Request) {
	var book Book
	parms := mux.Vars(r)

	//id, _ := strconv.Atoi(parms["id"])

	rows := db.QueryRow("Select * from books where id=$1", parms["id"])
	err := rows.Scan(&book.Id, &book.Title, &book.Author, &book.Year)
	logFatal(err)

	json.NewEncoder(w).Encode(book)
}

func addBook(w http.ResponseWriter, r *http.Request) {

	var book Book
	var bookID int

	json.NewDecoder(r.Body).Decode(&book)

	db.QueryRow("insert into books(title, author, year) values($1,$2,$3) RETURNING id;", book.Title, book.Author, book.Year).Scan(&bookID)
	//log.Fatal(err)

	json.NewEncoder(w).Encode(bookID)

}

func updateBook(w http.ResponseWriter, r *http.Request) {

	var book Book
	json.NewDecoder(r.Body).Decode(&book)

	db.Exec("update books set title=$1, author=$2,year=$3 where id=$4 RETURNING id", &book.Title, &book.Author, &book.Year, &book.Id)

	//log.Fatal(err)
	//rowUp, err := result.RowsAffected()

	json.NewEncoder(w).Encode(&book.Id)

}

func removeBook(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)

	result, err := db.Exec("delete from books where id = $1", params["id"])
	logFatal(err)

	rowsDeleted, err := result.RowsAffected()
	logFatal(err)

	json.NewEncoder(w).Encode(rowsDeleted)

}
