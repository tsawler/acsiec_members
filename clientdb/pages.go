package clientdb

import (
	"database/sql"
	"fmt"
	"github.com/tsawler/goblender/pkg/models"
)

// PageModel wraps database
type PageModel struct {
	DB *sql.DB
}

// AllPages returns slice of pages from goBlender's database
func (m *PageModel) AllPages() ([]*models.Page, error) {
	stmt := "SELECT id, page_title, active, slug, created_at, updated_at FROM pages ORDER BY page_title"

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pages []*models.Page

	for rows.Next() {
		var s models.Page
		err = rows.Scan(&s.ID, &s.PageTitle, &s.Active, &s.Slug, &s.CreatedAt, &s.UpdatedAt)
		if err != nil {
			return nil, err
		}
		pages = append(pages, &s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return pages, nil
}

func (m *PageModel) IsPrincipalPage(id int) bool {
	query := "SELECT count(page_id) from principal_pages where page_id = $1"
	row := m.DB.QueryRow(query, id)
	var num int
	err := row.Scan(&num)
	if err != nil {
		fmt.Println(err)
		return true
	}
	if num > 0 {
		return true
	}
	return false

}
