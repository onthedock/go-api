package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func DeletePerson(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{"message": "delete person with id " + id})
}
