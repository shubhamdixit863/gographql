package models

import (
	"github.com/saisri/gographql/internal/pkg/db"
	"log"
)

type Book struct {
	ID    string
	Title string
	User  User
}

func (book *Book) Save() {
	//#3
	stmt, err := db.Db.Prepare("INSERT INTO books(title,user_id) VALUES($1,$2)")
	if err != nil {
		log.Fatal(err)
	}
	//#4
	_, err = stmt.Exec(book.Title, book.User.ID)
	if err != nil {
		log.Fatal(err)
	}

	log.Print("Row inserted!")

}

func (book *Book) GetAll() []Book {
	stmt, err := db.Db.Prepare("select b.id, b.title, b.user_id, U.Username from books b inner join Users U on b.user_id = U.id")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var books []Book
	var username string
	var id string
	for rows.Next() {
		var book Book
		err := rows.Scan(&book.ID, &book.Title, &id, &username)
		if err != nil {
			log.Fatal(err)
		}
		book.User = User{
			ID:       id,
			Username: username,
		}
		books = append(books, book)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
	return books
}
