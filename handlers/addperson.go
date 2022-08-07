package handlers

import (
	"net/http"
	"personweb/models"

	"github.com/gin-gonic/gin"
)

func AddPerson(c *gin.Context) {
	// json is an instance of the models.Person struct
	var json models.Person
	// We try to fit the JSON object provided by the user
	// into the models.Person variable
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// If the "json" fits the models.Person struct, we try to insert it
	// into the database
	newId, err := models.AddPerson(json)
	// This version of the models.AddPerson function return the newID of
	// the record inserted (the old version just returned a bool value)
	// When there's a problem inserting the record, we return 0
	if newId != 0 {
		c.JSON(http.StatusOK, gin.H{"message": newId})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}
