package models

import (
	"database/sql"
	"log"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

// DB is a struct defined in the database/sql package
// DB stores connection information
var DB *sql.DB

func ConnectDatabase() error {
	// Try to open the database

	db, err := sql.Open("sqlite3", "./names.db")
	if err != nil {
		// Something went wrong
		// Return the err (and don't panic) so it can be handled elsewhere
		log.Print("[error] error opening database: " + err.Error())
		return err
	}

	// Everything Ok.
	// Store connection information in the global variable DB
	DB = db
	return nil
}

type Person struct {
	Id        int    `json:"id"`
	FirstName string `json:"firs_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	IPAddress string `json:"ip_address"`
}

func GetPersons(count int) ([]Person, error) {
	// Queries the DB
	rows, err := DB.Query("SELECT id, first_name, last_name, email, ip_address FROM people LIMIT " + strconv.Itoa(count))

	if err != nil {
		// Something went wrong
		log.Print("[error] DB.Query error " + err.Error())
		return nil, err
	}
	// We have a conection pointing to a bunch of rows
	// Use defer to make sure we close the connection.
	defer rows.Close()

	// Create a temporary, empty slice of Person
	// We'll add the results retrieved from the database to the slice
	people := make([]Person, 0)

	// Iterate over the rows (while there are any left)
	for rows.Next() {
		// Create an empty Person
		onePerson := Person{}
		// Retrieve the values from a row in the database and save them in the
		// designated fields of the struct
		// This action may fail (that's what we store returned value,
		// to check if it was an error)
		err = rows.Scan(&onePerson.Id, &onePerson.FirstName, &onePerson.LastName, &onePerson.Email, &onePerson.IPAddress)
		if err != nil {
			log.Print("[error] error scaning row " + err.Error())
			return nil, err
		}
		// We successfully retrieved the values from the row and store them
		// in the temporary struct. We add it to the slice of results.
		people = append(people, onePerson)
	}
	// rows.Next() return an error if something went wrong getting the next row.
	// So once the rows.Next() returns false, it is because there are no more rows
	// or because there was an error getting the next row?
	err = rows.Err()
	if err != nil {
		// There was an error getting the next row
		log.Print("[error] error getting next row: " + err.Error())
		return nil, err
	}
	// No more rows left, so we have retrieve everything we wanted
	return people, nil
}

func GetPersonById(id string, count int) ([]Person, error) {
	// Convert the number of records to retrieve (int) to string
	// to use it in the QUERY select.
	var strCount string = strconv.Itoa(count)
	// Retrieve the requested record (strCount defaults to 1)
	// If strCount != 1, we retrieve the strCount next records
	rows, err := DB.Query("SELECT id, first_name, last_name, email, ip_address FROM people WHERE id >= ? LIMIT ?", id, strCount)
	if err != nil {
		log.Printf("[error] error preparing query: %s", err.Error())
		return nil, err
	}
	// If there's no error, we make sure to close the connecttion
	defer rows.Close()

	// Here we will store what we retrieve from the DB
	results := make([]Person, 0)

	// Cycle through the results until there are no more records
	// or we get an error
	for rows.Next() {
		singlePerson := Person{}
		err = rows.Scan(&singlePerson.Id, &singlePerson.FirstName, &singlePerson.LastName, &singlePerson.Email, &singlePerson.IPAddress)
		if err != nil {
			log.Printf("[error] error scanning row: %s", err.Error())
			return nil, err
		}
		results = append(results, singlePerson)
	}
	if err != nil {
		log.Printf("[error] getting next row: %s", err.Error())
		return nil, err
	}

	return results, nil
}
