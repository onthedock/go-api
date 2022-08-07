package handlers

import (
	"net/http"
	"personweb/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func UpdatePerson(c *gin.Context) {
	var json models.Person
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pId, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	success, err := models.UpdatePerson(json, pId)
	if success {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
		return
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
}
