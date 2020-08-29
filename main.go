package main

import (
	"github.com/Pasidh/gofinal/middleware"
	"github.com/Pasidh/gofinal/task"
	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(middleware.Auth)

	r.GET("/customers", task.GetTodosHandler)
	r.POST("/customers", task.CreateTodosHandler)
	r.GET("/customers/:id", task.GetTodoByIdHandler)
	r.PUT("/customers/:id", task.UpdateTodosHandler)
	r.DELETE("/customers/:id", task.DeleteTodosHandler)

	return r
}

func main() {
	task.Init()
	task.CreateTable()
	r := setupRouter()
	r.Run(":2009")
}
