package models

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"
)

type BasicAuthor struct {
	Username string
	Email    string
	Password string
	Name     string
}

type Author struct {
	Id          int64
	CreatedTime time.Time
	BasicAuthor
}

func CreateAuthorTable(db *sql.DB) (sql.Result, error) {
	return db.Exec(`
    CREATE TABLE IF NOT EXISTS author (
      id SERIAL PRIMARY KEY ,
      username text UNIQUE NOT NULL ,
      email text UNIQUE NOT NULL ,
      password text NOT NULL ,
      name text NOT NULL ,
      created_time timestamp NOT NULL
    )
  `)
}

func AddOneAuthor(db *sql.DB, a BasicAuthor) {
	allAuthors := GetAllAuthors(db)
	for _, author := range allAuthors {
		if author.Username == a.Username || author.Email == a.Email {
			return
		}
	}
	_, err := db.Exec(
		"INSERT INTO author (username, email, password, name, created_time) VALUES ($1, $2, $3, $4, $5)",
		a.Username, a.Email, a.Password, a.Name, a.CreatedTime,
	)
	if err != nil {
		log.Fatalln(err)
	}
}

func GetAllAuthors(db *sql.DB) (authors []Author) {
	rows, err := db.Query("SELECT id, username, email, password, name, created_time FROM author")
	if err != nil {
		log.Fatalln(err)
	}
	for rows.Next() {
		var id int64
		var username string
		var email string
		var password string
		var name string
		var createdTime time.Time
		rows.Scan(&id, &username, &email, &password, &name, &createdTime)
		a := BasicAuthor{
			Username: username,
			Email:    email,
			Password: password,
			Name:     name,
		}
		authors = append(authors, Author{
			Id:          id,
			CreatedTime: createdTime,
			BasicAuthor: a,
		})
	}
	return authors
}

func GetOneAuthor(db *sql.DB, condition ...string) Author {
	query := "SELECT id, username, email, password, name, created_time FROM author"
	if len(condition) > 0 {
		query += " " + condition[0]
	}
	row := db.QueryRow(query)
	var id int64
	var username string
	var email string
	var password string
	var name string
	var createdTime time.Time
	row.Scan(&id, &username, &email, &password, &name, &createdTime)
	basicAuthor := BasicAuthor{
		Username: username,
		Email:    email,
		Password: password,
		Name:     name,
	}
	return Author{
		Id:          id,
		CreatedTime: createdTime,
		BasicAuthor: basicAuthor,
	}
}

func GetAuthorById(db *sql.DB, id int64) Author {
	condition := "WHERE id = " + string(id)
	return GetOneAuthor(db, condition)
}

func GetAuthorByEmail(db *sql.DB, email string) Author {
	condition := "WHERE email = " + email
	return GetOneAuthor(db, condition)
}

func UpdateOneAuthor(db *sql.DB, id int64, a BasicAuthor) (Author, error) {
	authors := GetAllAuthors(db)
	for _, author := range authors {
		if author.Username == a.Username {
			return Author{}, errors.New(fmt.Sprintf("Update author info failed, the username %s is duplicated.", a.Username))
		} else if author.Email == a.Email {
			return Author{}, errors.New(fmt.Sprintf("Update author info failed, the email %s is duplicated.", a.Email))
		}
	}
	query := fmt.Sprintf(
		"UPDATE author SET username = %s, email = %s， password = %s, name = %s WHERE id = %d",
		a.Username, a.Email, a.Password, a.Name, id,
	)
	_, err := db.Exec(query)
	if err != nil {
		log.Fatalln(err)
	}
	return GetAuthorById(db, id), nil
}

func DeleteOneAuthor(db *sql.DB, id int64) (Author, error) {
	author := GetAuthorById(db, id)
	query := fmt.Sprintf("DELETE FROM author WHERE id = %d", id)
	_, err := db.Exec(query)
	if err != nil {
		log.Fatalln(err)
	}
	return author, nil
}
