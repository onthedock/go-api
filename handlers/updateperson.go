package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func UpdatePerson(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{"message": "update person with id " + id})
}
