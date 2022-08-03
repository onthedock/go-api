package main

import (
	"log"
	"net/http"

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
	var count int = 10
	persons, err := models.GetPersons(count)
	if err != nil {
		log.Printf("[error] %v", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": "error"})
	}
	if persons == nil {
		c.JSON(http.StatusOK, gin.H{"message": "no records found"})
		return
	} else {
		log.Printf("[info] returned %v", persons)
		c.JSON(http.StatusOK, gin.H{"message": persons})
	}
}

func getPersonById(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{"message": "return person with id " + id})
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
