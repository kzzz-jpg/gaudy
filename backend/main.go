package main

import (
	"fmt"
	"guadb/repository"
)

func main() {
	db, err := repository.initdb()
	if err != nil {
		fmt.Println(err)
		return
	}
	guaRepo := repository.NewGuaRepo(db)

}
