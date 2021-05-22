package main

import "fmt"

type Book struct {
	Id     int
	Title  string
	Author string
	Year   string
}

var books []Book

func main() {

	books = append(books, Book{Id: 10, Title: "Golang ", Author: "ABC", Year: "2010"})
	var result []Book

	fmt.Println(books)

	var bookUpdate Book = Book{Id: 10, Title: "c++"}

	for _, item := range books {
		if item.Id == bookUpdate.Id {

			item.Title = bookUpdate.Title

		}
		result = append(result, item)
	}
	fmt.Println(result)
}
