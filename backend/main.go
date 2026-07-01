package main

import (
	"fmt"
	"guadb/handler"
	"guadb/repository"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	//如果環境沒有 zhparser 就用 repository.InitPostgreDB()
	db, err := repository.InitPostgreDBWithZhparser()
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

	fmt.Println("server is run at localhost:8080")
	if err := r.Run(":8080"); err != nil {
		fmt.Println(err)
		return
	}
}
