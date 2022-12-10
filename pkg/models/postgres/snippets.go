package postgres

import (
	"database/sql"
	"errors"
	"viktorkrams/snippetbox/pkg/models"
)

type SnippetModel struct {
	DB *sql.DB
}

func (model *SnippetModel) Insert(title, content string) (int, error) {
	transaction, err := model.DB.Begin()
	if err != nil {
		return 0, nil
	}
	query := `INSERT INTO snippets (title, content, created, expires)
	VALUES($1, $2, current_timestamp, current_timestamp + interval '365' day)
	RETURNING id`
	lastInsertedId := 0

	err = transaction.QueryRow(query, title, content).Scan(&lastInsertedId)
	if err != nil {
		transaction.Rollback()
		return 0, err
	}

	err = transaction.Commit()
	if err != nil {
		return 0, err
	}

	return lastInsertedId, nil
}

func (model *SnippetModel) Get(id int) (*models.Snippet, error) {
	transaction, err := model.DB.Begin()
	if err != nil {
		return nil, err
	}
	query := "SELECT * FROM snippets WHERE expires > current_timestamp AND id = $1"
	snippet := &models.Snippet{}
	err = transaction.QueryRow(query, id).Scan(&snippet.ID, &snippet.Title, &snippet.Content, &snippet.Created, &snippet.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}

		return nil, err
	}

	err = transaction.Commit()
	if err != nil {
		return nil, err
	}

	return snippet, nil
}

func (model *SnippetModel) Latest() ([]*models.Snippet, error) {
	query := "SELECT * FROM snippets WHERE expires > current_timestamp ORDER BY created DESC LIMIT 10"
	rows, err := model.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var snippets []*models.Snippet

	for rows.Next() {
		snippet := &models.Snippet{}
		err := rows.Scan(&snippet.ID, &snippet.Title, &snippet.Content, &snippet.Created, &snippet.Expires)
		if err != nil {
			return nil, err
		}

		snippets = append(snippets, snippet)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}

func (model *SnippetModel) Delete(id int) error {
	transaction, err := model.DB.Begin()
	if err != nil {
		return err
	}

	query := "DELETE FROM snippets WHERE id = $1"
	transaction.QueryRow(query, id)

	err = transaction.Commit()
	if err != nil {
		return err
	}

	return nil
}
