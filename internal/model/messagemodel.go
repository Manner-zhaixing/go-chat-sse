package model

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ MessageModel = (*customMessageModel)(nil)

type (
	// MessageModel is an interface to be customized, add more methods here,
	// and implement the added methods in customMessageModel.
	MessageModel interface {
		messageModel
		FindOneBySessionId(ctx context.Context, sessionId int64) (*Message, error)
	}

	customMessageModel struct {
		*defaultMessageModel
	}
)

// NewMessageModel returns a model for the database table.
func NewMessageModel(conn sqlx.SqlConn) MessageModel {
	return &customMessageModel{
		defaultMessageModel: newMessageModel(conn),
	}
}

// FindOneBySessionId 通过session获取message
func (c customMessageModel) FindOneBySessionId(ctx context.Context, sessionId int64) (*Message, error) {
	var resp Message
	query := fmt.Sprintf("select %s from %s where `session_id` = ? limit 1", messageRows, c.table)
	err := c.conn.QueryRowCtx(ctx, &resp, query, sessionId)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}
