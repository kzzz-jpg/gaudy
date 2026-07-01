package repository

import (
	"database/sql"
	"guadb/model"
	"strings"

	"github.com/lib/pq"
)

type GuaRepo interface {
	AddGua(*model.Gua) (int, error)
	GetGua(*model.Gua) ([]*model.Gua, error)
}

// 如果環境沒有 zhparser 就 return &guaRepo{db}
func NewGuaRepo(db *sql.DB) GuaRepo {
	return &guaRepoWithZh{db}
}

type guaRepo struct {
	db *sql.DB
}

func (g *guaRepo) AddGua(gua *model.Gua) (int, error) {
	var guaid int
	err := g.db.QueryRow(`
    	INSERT INTO guas (title,people,people_str,content)VALUES($1,$2,$3,$4)
    	RETURNING gua_id;
	`, gua.Title, pq.Array(gua.People), strings.Join(gua.People, " "), gua.Content).Scan(&guaid)
	if err != nil {
		return 0, err
	}
	return guaid, nil
}
func (g *guaRepo) GetGua(gua *model.Gua) ([]*model.Gua, error) {
	var guas []*model.Gua

	parts := []string{}
	if gua.Title != "" {
		parts = append(parts, gua.Title)
	}
	parts = append(parts, gua.People...)
	if gua.Content != "" {
		parts = append(parts, gua.Content)
	}
	searchStr := strings.Join(parts, " or ")

	row, err := g.db.Query(`
		SELECT gua_id,title,people,content FROM guas 
		    WHERE search_vector @@ websearch_to_tsquery('simple',$1)
		ORDER BY ts_rank(search_vector,websearch_to_tsquery('simple',$1)) DESC
    `, searchStr)
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
