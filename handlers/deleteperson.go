package handlers

import (
	"net/http"
	"personweb/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func DeletePerson(c *gin.Context) {
	var id int = 0
	var err error
	id, err = strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	success, err := models.DeletePerson(id)
	if success {
		c.JSON(http.StatusOK, gin.H{"message": "deleted record"})
		return
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
}
