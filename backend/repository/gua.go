package repository

import (
	"database/sql"
	"guadb/model"
)

type GuaRepo interface {
	AddGua(*model.Gua) (int, error)
	GetGua(*model.Gua) ([]*model.Gua, error)
}

type guaRepo struct {
	db *sql.DB
}

func NewGuaRepo(db *sql.DB) GuaRepo {
	return &guaRepo{db}
}

func (g *guaRepo) AddGua(gua *model.Gua) (int, error) {
	var guaid int
	err := g.db.QueryRow(`
    	INSERT INTO guas (title,people,content)VALUES($1,$2,$3)
    	RETURNING gua_id;
	`, gua.Title, gua.People, gua.Content).Scan(&guaid)
	if err != nil {
		return -1, err
	}
	return guaid, nil
}
func (g *guaRepo) GetGua(gua *model.Gua) ([]*model.Gua, error) {
	var guas []*model.Gua
	row, err := g.db.Query(`
		SELECT gua_id,title,people,content FROM guas 
		    WHERE title ILIKE '%' || REPLACE(REPLACE($1,'%','\%'),'_','\_') || '%' ESCAPE '\'
		    OR content ILIKE '%' || REPLACE(REPLACE($2,'%','\%'),'_','\_') || '%' ESCAPE '\'
			OR EXISTS(
			    SELECT 67 FROM unnest(people) AS P WHERE P ILIKE '%' || REPLACE(REPLACE($3,'_','\_'),'%','\%') || '%'
			)
    `, gua.Title, gua.People, gua.Content)
	if err != nil {
		return nil, err
	}
	for row.Next() {
		var gua model.Gua
		err = row.Scan(&gua.Title, &gua.People, &gua.Content)
		if err != nil {
			return nil, err
		}
		guas = append(guas, &gua)
	}
	if err != nil {
		return nil, err
	}
	return guas, nil
}
