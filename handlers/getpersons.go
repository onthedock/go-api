package handlers

import (
	"log"
	"net/http"
	"personweb/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetPersons(c *gin.Context) {
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
