package main

import (
	"log"
	"personweb/handlers"
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
		v1.GET("person", handlers.GetPersons)
		v1.GET("person/:id", handlers.GetPersonById)
		v1.POST("person", handlers.AddPerson)
		v1.PUT("person/:id", handlers.UpdatePerson)
		v1.DELETE("person/:id", handlers.DeletePerson)
	}
	r.Run()
}
