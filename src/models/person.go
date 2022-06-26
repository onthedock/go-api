package models

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func ConnectDatabase() error {
	db, err := sql.Open("sqlite3", "./names.db")
	if err != nil {
		//log.Printf("Error opening database ./names.db: %s\n", err.Error())
		return err
	}
	DB = db
	return nil
}

type Person struct {
	Id        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	IpAddress string `json:"ip_address"`
}

func GetPersons(count int) ([]Person, error) {
	rows, err := DB.Query("SELECT id, first_name, last_name, email, ip_address from people LIMIT " + strconv.Itoa(count))

	if err != nil {
		log.Printf("Error querying the database: %s\n", err.Error())
		return nil, err
	}
	defer rows.Close()

	people := make([]Person, 0)

	for rows.Next() {
		singlePerson := Person{}
		err = rows.Scan(&singlePerson.Id, &singlePerson.FirstName, &singlePerson.LastName, &singlePerson.Email, &singlePerson.IpAddress)

		if err != nil {
			log.Printf("Error getting the record with Scan: %s", err.Error())
			return nil, err
		}

		people = append(people, singlePerson)
	}

	err = rows.Err()
	if err != nil {
		log.Printf("Error getting the Next: %s", err.Error())
		return nil, err
	}

	return people, err
}

func GetPersonById(id string) (Person, error) {

	stmt, err := DB.Prepare("SELECT id, first_name, last_name, email, ip_address from people WHERE id = ?")

	if err != nil {
		log.Printf("Error querying the database: %s", err.Error())
		return Person{}, err
	}

	person := Person{}

	sqlErr := stmt.QueryRow(id).Scan(&person.Id, &person.FirstName, &person.LastName, &person.Email, &person.IpAddress)

	if sqlErr != nil {
		if sqlErr == sql.ErrNoRows {
			return Person{}, nil
		}
		return Person{}, sqlErr
	}
	return person, nil
}

func AddPerson(newPerson Person) (bool, error) {

	tx, err := DB.Begin()
	if err != nil {
		log.Printf("Error opening the database: %s", err.Error())
		return false, err
	}
	stmt, err := tx.Prepare("INSERT INTO people (first_name, last_name, email, ip_address) VALUES (?,?,?,?)")
	if err != nil {
		log.Printf("error preparing insert statament: %s", err.Error())
		return false, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(newPerson.FirstName, newPerson.LastName, newPerson.Email, newPerson.IpAddress)
	if err != nil {
		log.Printf("error executing the query: %s", err.Error())
		return false, err
	}
	tx.Commit()
	return true, nil
}

func UpdatePerson(ourPerson Person, id int) (bool, error) {
	tx, err := DB.Begin()
	if err != nil {
		fmt.Printf("Error connecting to db, %s", err.Error())
		return false, err
	}
	stmt, err := tx.Prepare("UPDATE people SET first_name = ?, last_name = ?, email = ?, ip_address = ? WHERE Id = ?")
	if err != nil {
		fmt.Printf("Error preparing update: %s", err.Error())
		return false, err
	}
	defer stmt.Close()
	_, err = stmt.Exec(ourPerson.FirstName, ourPerson.LastName, ourPerson.Email, ourPerson.IpAddress, id)
	if err != nil {
		fmt.Printf("Error executing the query: %s", err.Error())
		return false, err
	}
	tx.Commit()
	return true, nil
}

func DeletePerson(personId int) (bool, error) {
	tx, err := DB.Begin()
	if err != nil {
		fmt.Printf("Error connecting to database: %s", err.Error())
		return false, err
	}
	stmt, err := DB.Prepare("DELETE FROM people WHERE id = ?")
	if err != nil {
		fmt.Printf("Error preparing statement: %s", err.Error())
		return false, err
	}
	defer stmt.Close()
	_, err = stmt.Exec(personId)
	if err != nil {
		fmt.Printf("Error deleting record: %s", err.Error())
		return false, err
	}
	tx.Commit()
	return true, nil
}
