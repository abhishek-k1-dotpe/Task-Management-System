package models

import (
	"context"
	"database/sql"
	"errors"
	"user-service/db"
	"user-service/domain"

	"github.com/Masterminds/squirrel"
)

func Create(ctx context.Context, user domain.User) (int, error) {

	var conn *sql.Conn
	conn, err := db.Client.Conn(ctx)
	if err != nil {
		return 0, err
	}
	defer conn.Close()
	val := []interface{}{user.Username, user.Email, user.Type, user.CreatedBy}
	qbuilder := squirrel.Insert("users").Columns("name", "email", "type", "created_by")
	qbuilder = qbuilder.Values(val...)
	query, _, err := qbuilder.ToSql()
	if err != nil {
		return 0, err
	}
	result, err := conn.ExecContext(ctx, query, val...)

	if err != nil {
		return 0, err
	}
	userId, _ := result.LastInsertId()
	return int(userId), err
}

func DeleteUser(ctx context.Context, userId int) error {
	var conn *sql.Conn
	conn, err := db.Client.Conn(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()
	qbuilder := squirrel.Delete("users").Where(squirrel.Eq{"user_id": userId})
	query, qargs, err := qbuilder.ToSql()
	if err != nil {
		return err
	}
	_, err = conn.ExecContext(ctx, query, qargs...)
	if err != nil {
		return err
	}
	return nil

}

func GetUserInfo(ctx context.Context, userId int) (*domain.User, error) {
	var conn *sql.Conn
	conn, err := db.Client.Conn(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	user := domain.User{}
	qbuilder := squirrel.Select("user_id", "name", "email", "created_by", "type").From("users").Where(squirrel.Eq{"user_id": userId})
	query, qargs, err := qbuilder.ToSql()
	if err != nil {
		return nil, err
	}
	row := conn.QueryRowContext(ctx, query, qargs...)
	err = row.Scan(&user.UserId, &user.Username, &user.Email, &user.CreatedBy, &user.Type)

	if err == sql.ErrNoRows {
		return nil, errors.New("no such user_id in the system")
	} else if err != nil {
		return nil, err
	}
	return &user, nil
}
