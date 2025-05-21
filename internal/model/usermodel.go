package model

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ UserModel = (*customUserModel)(nil)

type (
	// UserModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserModel.
	UserModel interface {
		userModel
		FindByUsername(ctx context.Context, username string) (*User, error)
		FindByUsernameAndPassword(ctx context.Context, username string, password string) (*User, error)
		UpdateConversationNumByConversationIdAndUserId(ctx context.Context, UserId int64) error
	}

	customUserModel struct {
		*defaultUserModel
	}
)

// NewUserModel returns a model for the database table.
func NewUserModel(conn sqlx.SqlConn) UserModel {
	return &customUserModel{
		defaultUserModel: newUserModel(conn),
	}
}

func (m *customUserModel) FindByUsername(ctx context.Context, username string) (*User, error) {
	//TODO implement me
	query := fmt.Sprintf("select %s from %s where `username` = ? limit 1", userRows, m.table)
	var resp User
	err := m.conn.QueryRowCtx(ctx, &resp, query, username)
	switch err {
	case nil:
		return &resp, nil
	case sql.ErrNoRows, sqlx.ErrNotFound:
		return nil, nil
	default:
		return nil, err
	}
}

func (m *customUserModel) FindByUsernameAndPassword(ctx context.Context, username string, password string) (*User, error) {
	//TODO implement me
	query := fmt.Sprintf("select %s from %s where `username` = ? and `password` = ? limit 1", userRows, m.table)
	var resp User
	err := m.conn.QueryRowCtx(ctx, &resp, query, username, password)
	switch err {
	case nil:
		return &resp, nil
	case sql.ErrNoRows, sqlx.ErrNotFound:
		return nil, nil
	default:
		return nil, err
	}
}

func (m *customUserModel) UpdateConversationNumByConversationIdAndUserId(ctx context.Context, userID int64) error {
	const updateQuery = "UPDATE %s SET conversation_nums = conversation_nums + 1 WHERE id = ?"
	_, err := m.conn.ExecCtx(ctx, fmt.Sprintf(updateQuery, m.table), userID)
	if err != nil {
		return fmt.Errorf("failed to increment conversation_nums: %w", err)
	}

	return nil
}
