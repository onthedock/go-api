package handlers

import (
	"log"
	"net/http"
	"personweb/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetPersonById(c *gin.Context) {
	// Limit the number of records returned
	var max_count int = 10
	// Default value
	var count int = 1

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
