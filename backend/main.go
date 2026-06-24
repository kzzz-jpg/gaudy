package main

import (
	"fmt"
	"guadb/handler"
	"guadb/repository"

	"github.com/gin-gonic/gin"
)

func main() {
	db, err := repository.Initdb()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()
	guaRepo := repository.NewGuaRepo(db)
	guaHandler := handler.NewGuaHandler(guaRepo)

	r := gin.Default()
	api := r.Group("/api")
	{
		api.POST("/gua", guaHandler.AddGuaHandler)
		api.POST("/gua/list", guaHandler.GetGuaHandler)
	}

	if err := r.Run(":8080"); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("server is run at localhost:8080")
}
