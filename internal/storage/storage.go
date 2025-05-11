package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/philipslstwoyears/calculator-go/internal/model/dto"
	"sort"
)

type Storage interface {
	AddExpression(e dto.Expression) (int, error)
	GetExpression(id int) (dto.Expression, bool)
	UpdateExpression(e dto.Expression) error
	GetExpressions(userID int) ([]dto.Expression, error)
	AddUser(e dto.User) (int, error)
	GetUser(login string) (dto.User, bool)
}

type DbStorage struct {
	db *sql.DB
}

func New(db *sql.DB) *DbStorage {
	return &DbStorage{
		db: db,
	}
}
func (s *DbStorage) AddExpression(e dto.Expression) (int, error) {
	var q = `
	INSERT INTO expressions (expression, user_id, result, status) values ($1, $2, $3, $4)
	`
	result, err := s.db.Exec(q, e.Expression, e.UserID, e.Result, e.Status)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (s *DbStorage) UpdateExpression(e dto.Expression) error {
	q := `
	UPDATE expressions
	SET expression = ?, user_id = ?, result = ?, status = ?
	WHERE id = ?
	`
	_, err := s.db.Exec(q, e.Expression, e.UserID, e.Result, e.Status, e.ID)
	return err
}

func (s *DbStorage) GetExpression(id int) (dto.Expression, bool) {
	var e dto.Expression
	q := `
	SELECT id, expression, user_id, result, status
	FROM expressions
	WHERE id = ?
	`
	err := s.db.QueryRow(q, id).Scan(&e.ID, &e.Expression, &e.UserID, &e.Result, &e.Status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return dto.Expression{}, false
		}
		return dto.Expression{}, false
	}
	return e, true
}

func (s *DbStorage) GetExpressions(userID int) ([]dto.Expression, error) {
	var expressions []dto.Expression

	q := `
	SELECT id, expression, user_id, result, status
	FROM expressions
	WHERE user_id = ?
	`

	rows, err := s.db.Query(q, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var e dto.Expression
		err := rows.Scan(&e.ID, &e.Expression, &e.UserID, &e.Result, &e.Status)
		if err != nil {
			return nil, err
		}
		expressions = append(expressions, e)
	}
	sort.Slice(expressions, func(i, j int) bool {
		return expressions[i].ID < expressions[j].ID
	})
	return expressions, nil
}

func (s *DbStorage) AddUser(e dto.User) (int, error) {
	// Используем плейсхолдеры SQLite
	query := `INSERT INTO users (login, password) VALUES (?, ?)`

	result, err := s.db.Exec(query, e.Login, e.Password)
	if err != nil {
		return 0, fmt.Errorf("AddUser insert error: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("AddUser lastInsertId error: %w", err)
	}

	return int(id), nil
}

func (s *DbStorage) GetUser(login string) (dto.User, bool) {
	var e dto.User

	// Выведем лог для отладки
	fmt.Println("Trying to get user with login:", login)

	query := `SELECT id, login, password FROM users WHERE login = ?`

	err := s.db.QueryRow(query, login).Scan(&e.Id, &e.Login, &e.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			fmt.Println("User not found in database")
			return dto.User{}, false
		}
		fmt.Println("GetUser query error:", err)
		return dto.User{}, false
	}

	return e, true
}
