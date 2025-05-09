package storage

import (
	"database/sql"
	"errors"
	"github.com/philipslstwoyears/calculator-go/internal/dto"
	"sort"
)

type Storage struct {
	db *sql.DB
}

func New(db *sql.DB) *Storage {
	return &Storage{
		db: db,
	}
}
func (s *Storage) AddExpression(e dto.Expression) (int, error) {
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

func (s *Storage) UpdateExpression(e dto.Expression) error {
	q := `
	UPDATE expressions
	SET expression = ?, user_id = ?, result = ?, status = ?
	WHERE id = ?
	`
	_, err := s.db.Exec(q, e.Expression, e.UserID, e.Result, e.Status, e.ID)
	return err
}

func (s *Storage) GetExpression(id int) (dto.Expression, bool) {
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

func (s *Storage) GetExpressions(userID int) []dto.Expression {
	var expressions []dto.Expression

	q := `
	SELECT id, expression, user_id, result, status
	FROM expressions
	WHERE user_id = ?
	`

	rows, err := s.db.Query(q, userID)
	if err != nil {
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		var e dto.Expression
		err := rows.Scan(&e.ID, &e.Expression, &e.UserID, &e.Result, &e.Status)
		if err != nil {
			continue
		}
		expressions = append(expressions, e)
	}
	sort.Slice(expressions, func(i, j int) bool {
		return expressions[i].ID < expressions[j].ID
	})
	return expressions
}

func (s *Storage) AddUser(e dto.User) (int, error) {
	var q = `
	INSERT INTO users (login, password) values ($1, $2)
	`
	result, err := s.db.Exec(q, e.Login, e.Password)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (s *Storage) GetUser(login string) (dto.User, bool) {
	var e dto.User
	q := `
	SELECT password, user_id, login
	FROM users
	WHERE login = ?
	`
	err := s.db.QueryRow(q, login).Scan(&e.Password, &e.Id, &e.Login)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return dto.User{}, false
		}
		return dto.User{}, false
	}
	return e, true
}
