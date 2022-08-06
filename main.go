package main

import (
	"log"
	"net/http"
	"strconv"

	"personweb/models"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	err := models.ConnectDatabase()
	if err != nil {
		log.Print("[error] error connecting to database " + err.Error())
	}
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

func getPersons(c *gin.Context) {
	// Limit the number of records returned
	var max_count int = 10
	// Default value
	var count int = 10

	// If no "count" parameter is provided, we return the count value
	if c.Query("count") != "" {
		var conv_err error
		count, conv_err = strconv.Atoi(c.Query("count"))
		if conv_err != nil {
			log.Printf("[error] error getting count from queryString (default to %d). error: %s", count, conv_err.Error())
		}
	}

	// Return max_count (at most)
	if count > max_count {
		log.Printf("[warning] requested %d records (returning max: %d)", count, max_count)
		count = max_count
	}

	persons, err := models.GetPersons(count)

	if err != nil {
		log.Printf("[error] %v", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": "error"})
	}

	if persons == nil {
		c.JSON(http.StatusOK, gin.H{"message": "no records found"})
		return
	} else {
		// log.Printf("[info] returned %v", persons)
		c.JSON(http.StatusOK, gin.H{"message": persons})
	}
}

func getPersonById(c *gin.Context) {
	// Limit the number of records returned
	var max_count int = 10
	// Default value
	var count int = 10

	// If no "count" parameter is provided, we return the count value
	if c.Query("count") != "" {
		var conv_err error
		count, conv_err = strconv.Atoi(c.Query("count"))
		if conv_err != nil {
			log.Printf("[error] error getting count from queryString (default to %d). error: %s", count, conv_err.Error())
			return
		}
	}

	// Return max_count (at most)
	if count > max_count {
		log.Printf("[warning] requested %d records (returning max: %d)", count, max_count)
		count = max_count
	}
	id := c.Param("id")

	results, err := models.GetPersonById(id, count)
	if err != nil {
		log.Printf("[error] retrieving record from database: %s", err.Error())
	}
	if results[0].FirstName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "not records found"})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"data": results})
	}

}

func addPerson(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "add person"})
}

func updatePerson(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{"message": "update person with id " + id})
}

func deletePerson(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{"message": "delete person with id " + id})
}
