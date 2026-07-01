package repository

import (
	"database/sql"
	"fmt"
	"guadb/model"
	"strings"

	"github.com/lib/pq"
)

type guaRepoWithZh struct {
	db *sql.DB
}

func (g *guaRepoWithZh) AddGua(gua *model.Gua) (int, error) {
	opsql := `
		INSERT INTO guas(title,people,people_str,content) VALUES($1,$2,$3,$4) RETURNING gua_id;
    `

	var id int
	if err := g.db.QueryRow(opsql, gua.Title, pq.Array(gua.People), strings.Join(gua.People, " "), gua.Content).Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}
func (g *guaRepoWithZh) GetGua(gua *model.Gua) ([]*model.Gua, error) {
	opsql := `
		SELECT gua_id,title,people,content,created_at FROM guas,websearch_to_tsquery('zhcfg',$1) AS q WHERE search_vector @@ q ORDER BY ts_rank(search_vector,q) DESC
    `
	opstr := fmt.Sprintf("%s %s %s", gua.Title, strings.Join(gua.People, " "), gua.Content)
	row, err := g.db.Query(opsql, opstr)
	if err != nil {
		return nil, err
	}
	defer row.Close()
	retgua := make([]*model.Gua, 0)
	for row.Next() {
		var tmpgua model.Gua
		err := row.Scan(&tmpgua.GuaId, &tmpgua.Title, pq.Array(&tmpgua.People), &tmpgua.Content, &tmpgua.CreatedAt)
		if err != nil {
			return nil, err
		}
		retgua = append(retgua, &tmpgua)
	}
	if row.Err() != nil {
		return nil, row.Err()
	}
	return retgua, nil
}
