package sqlpkg

import (
	"forum/model"
)

/*
inserts a new comment into DB, returns an ID for the comment
*/
func (f *ForumModel) GetCategories() ([]*model.Category, error) {
	q := `SELECT id, name FROM categories ORDER BY name`
	rows, err := f.DB.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// parsing the query's result
	var categories []*model.Category
	for rows.Next() {
		category:= &model.Category{}
		err = rows.Scan(&category.ID,&category.Name)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return categories, nil
}
