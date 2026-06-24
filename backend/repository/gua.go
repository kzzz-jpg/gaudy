package repository

import (
	"database/sql"
	"guadb/model"

	"github.com/lib/pq"
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
	`, gua.Title, pq.Array(gua.People), gua.Content).Scan(&guaid)
	if err != nil {
		return 0, err
	}
	return guaid, nil
}
func (g *guaRepo) GetGua(gua *model.Gua) ([]*model.Gua, error) {
	var guas []*model.Gua
	row, err := g.db.Query(`
		SELECT gua_id,title,people,content FROM guas 
		    WHERE ($1 <> '' AND title ILIKE '%' || REPLACE(REPLACE($1,'%','\%'),'_','\_') || '%' ESCAPE '\')
		    OR ($2 <> '' AND content ILIKE '%' || REPLACE(REPLACE($2,'%','\%'),'_','\_') || '%' ESCAPE '\')
			OR EXISTS(
			    SELECT 67 FROM unnest(people) P,unnest($3::TEXT[]) Q WHERE (Q <> '' AND P ILIKE '%' || REPLACE(REPLACE(Q,'_','\_'),'%','\%') || '%' ESCAPE '\')
			)
    `, gua.Title, gua.Content, pq.Array(gua.People))
	if err != nil {
		return nil, err
	}
	defer row.Close()
	for row.Next() {
		var item model.Gua
		err = row.Scan(&item.GuaId, &item.Title, pq.Array(&item.People), &item.Content)
		if err != nil {
			return nil, err
		}
		guas = append(guas, &item)
	}
	if err := row.Err(); err != nil {
		return nil, err
	}
	return guas, nil
}
