package repository

import (
	"database/sql"
	"guadb/model"
)

type guaRepoWithZh {
	db *sql.DB
}

func (g *guaRepoWithZh) AddGua(*model.Gua)(int, error){

}
func (g *guaRepoWithZh) GetGua(*model.Gua) (*model.Gua, error){

}

