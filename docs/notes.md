# API using Gin framework and SQLite backend

## Init module

```shell
go mod init personweb
```

## Go get dependecies

```shell
go get -u github.com/mattn-go-sqlite3
go get -u github.com/gin-gonic/gin
```

## Create the (empty) `main.go``

```shell
touch main.go
```

## Create the basic API (non-functional yet)

- Import the `gin` package.
- Define an *API group*
- Define the *API* (undefined) functions

```go
package main

import "github.com/gin-gonic/gin"

func main() {
 r := gin.Default()
 v1 := r.Group("/api/v1")
 {
  v1.GET("person", getPersons)
  v1.GET("person/:id", getPersonById)
  v1.POST("person", addPerson)
  v1.PUT("person/:id", updatePerson)
  v1.DELETE("person/:id", deletePerson)
 }
 r.Run()
}
```

These functions take care of the *incoming* request from the user; they check if the required parameters are present and so on before passing the request to the backend.

The *technical* name of these functions is *handlers*.

## Create *dummy* handler functions

We create *dummy* handler functions that simply return a message confirming that they are been successfully being called.

We learn how to get a parameter from the URI.

```go
// Example
func updatePerson(c *gin.Context) {
 id := c.Param("id")
 c.JSON(http.StatusOK, gin.H{"message": "update person with id " + id})
}
```

We also learn how to return JSON documents from using `gin.H{}`.

## Database

We dowload a prepopulated SQLite3 database from <https://github.com/JeremyMorgan/GoSqliteDemo/blob/main/names.db>. This database is part of the tutorial [How to use SQLite with Go](https://www.allhandsontech.com/programming/golang/how-to-use-sqlite-with-go/), from the same author.

## Create the folder `models`

We'll use the folder to store the `models` package.

We create the `models/person.go` file:

```shell
mkdir models
touch models/person.go
```

This model will host the functions to connect to the database.

We import the `database/sql` pckage (from the *standar* library) to interact with SQL databases. The package does not provide *drivers*, so we have to manually import the specific driver for the database that we use. In our case, SQLite.

To do so, we use the *blank identifier* (`_`), as we do not want to import any function from the package, just to *initialize* the driver so it is available to be used by the `database/sql` package.

```go
import (
 "database/sql"
 _ "github.com/mattn/go-sqlite3"
)
```

We define a global variable to *store* the database connection information; (`DB` is a *struct* defined in the `database/sql` package).

```go
var DB *sql.DB
```

Then, we define the database to connect to the database:

```go
func ConnectDatabase() error {
 db, err := sql.Open("sqlite3", "./names.db")
 if err != nil {
  return err
 }

 DB = db
 return nil
}
```

## Create the model

In the `models/person.go` file we define the *struct* that will hold our model. It will store values coming from the request, before we save them to the database or data from the database before we return it to the user.

We use JSON labels to define how we want the fields in the *struct* to be identified in the JSON document.

```go
type Person struct {
 Id        int    `json:"id"`
 FirstName string `json:"firs_name"`
 LastName  string `json:"last_name"`
 Email     string `json:"email"`
 IPAddress string `json:"ip_address"`
}
```

## Connect to the database

In the `main.go` file, we connect to the database:

```go
 err := models.ConnectDatabase()
 if err != nil {
  log.Print("[error] error connecting to database " + err.Error())
 }
```

## Get all (almost) *Person* records from the database

In the `models/person.go` we define the `GetPersons` function to retrieve all the records from the database.

> As the database may contain a huge number of records (the `names.db` database contains 1000 entries), we limit the number of results to 10.

We execute the query and check if there was some error. If no error happened, we have a "pointer" to a bunch of rows. We use `defer` to make sure we release the "pointer" when it's no longer needed.

We create an empty slice of `Person` to store (and return) the results retrieved from the database.

We iterate over the rows resulting from the query with `rows.Next()`. We create an empty `Person`; we will store the values in each field in the row with the field defined in the struct.

For each row, using `row.Scan()` we retrieve the fields in the row and store them in the struct's fields.

This may fail for whatever reason, so we check for errors in `err := rows.Scan()`.

If nothing failed, we've "copied" the values from the retrieved row into the fields of the `Person` struct. We add the the `Person` struc to the slice of `Person` where we save the results retrieved from the database.

`rows.Next()` may fail if there is a problem getting the *next* row. In that case, it returns `false`, so we get out of the loop of retrieving rows. But it also returns `false` when there are no more rows.

To check if there was an error getting the next row or we just run out of rows to retrieve, we get `err := row.Error()` and we check it. If there was no error, this means that we have retrieved all rows, so we can return the slice of `Person` with all the values (and no error).

## Calling the `models.GetPersons()` function

The `models.GetPersons()` retrieves results from the database, but it's not called (yet) from the `main.go` file.

We modify the `getPersons()` function to call the `models.GetPerson()` function to retrieve 10 records from the database.
