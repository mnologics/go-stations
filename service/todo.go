package service

import (
	"context"
	"database/sql"

	"github.com/TechBowl-japan/go-stations/model"
)

// A TODOService implements CRUD of TODO entities.
type TODOService struct {
	db *sql.DB
}

// NewTODOService returns new TODOService.
func NewTODOService(db *sql.DB) *TODOService {
	return &TODOService{
		db: db,
	}
}

// CreateTODO creates a TODO on DB.
func (s *TODOService) CreateTODO(ctx context.Context, subject, description string) (*model.TODO, error) {
	const (
		insert  = `INSERT INTO todos(subject, description) VALUES(?, ?)`
		confirm = `SELECT id, subject, description, created_at, updated_at FROM todos WHERE id = ?`
	)
	// log.Println("CreateTODO", subject, description)
	result, err := s.db.ExecContext(ctx, insert, subject, description)
	if err != nil {
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	var todo model.TODO
	err = s.db.QueryRowContext(ctx, confirm, id).Scan(&todo.ID, &todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt)
	if err != nil {
		return nil, err
	}
	// log.Println("todo", todo)
	return &todo, nil
}

// ReadTODO reads TODOs on DB.
func (s *TODOService) ReadTODO(ctx context.Context, prevID, size int64) ([]*model.TODO, error) {
	const (
		read       = `SELECT id, subject, description, created_at, updated_at FROM todos ORDER BY id ASC LIMIT ?`
		readWithID = `SELECT id, subject, description, created_at, updated_at FROM todos WHERE id < ? ORDER BY id ASC LIMIT ?`
	)
	// log.Println("ReadTODO", prevID, size)
	rows, err := s.db.QueryContext(ctx, readWithID, prevID, size)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// log.Println("rows", rows)

	var todos []*model.TODO
	for rows.Next() {
		var todo model.TODO
		if err := rows.Scan(&todo.ID, &todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt); err != nil {
			return nil, err
		}
		// log.Println("todo", todo)
		todos = append(todos, &todo)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	// log.Println("todos", todos)
	return todos, nil
}

// UpdateTODO updates the TODO on DB.
func (s *TODOService) UpdateTODO(ctx context.Context, id int64, subject, description string) (*model.TODO, error) {
	const (
		update  = `UPDATE todos SET subject = ?, description = ? WHERE id = ?`
		confirm = `SELECT id, subject, description, created_at, updated_at FROM todos WHERE id = ?`
	)
	// log.Println("UpdateTODO", id, subject, description)
	// log.Printf("id=%d, subject=%s, description=%s", id, subject, description)
	// log.Printf("UPDATE todos SET subject = %s, description = %s WHERE id = %d", subject, description, id)
	result, err := s.db.ExecContext(ctx, update, subject, description, id)
	if err != nil {
		// log.Println("err:", err)
		return nil, err // subjectが空文字の場合に対応するため (ステーション12)
	}
	_ = result
	// log.Println("result", result)

	var todo model.TODO
	err = s.db.QueryRowContext(ctx, confirm, id).Scan(&todo.ID, &todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt)
	if err != nil {
		// log.Println("err:", err)
		return nil, &model.ErrNotFound{} // idに対応するTODOが存在しない場合に対応するため (ステーション12)
	}
	// log.Println("todo", todo)
	return &todo, nil
}

// DeleteTODO deletes TODOs on DB by ids.
func (s *TODOService) DeleteTODO(ctx context.Context, ids []int64) error {
	const deleteFmt = `DELETE FROM todos WHERE id IN (?%s)`

	return nil
}
