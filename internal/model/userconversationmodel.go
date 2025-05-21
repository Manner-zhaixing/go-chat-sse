package model

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ UserConversationModel = (*customUserConversationModel)(nil)

type (
	// UserConversationModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserConversationModel.
	UserConversationModel interface {
		userConversationModel
		DeleteByConversationIdAndUserId(ctx context.Context, conversationId int64, userId int64) error
	}

	customUserConversationModel struct {
		*defaultUserConversationModel
	}
)

// NewUserConversationModel returns a model for the database table.
func NewUserConversationModel(conn sqlx.SqlConn) UserConversationModel {
	return &customUserConversationModel{
		defaultUserConversationModel: newUserConversationModel(conn),
	}
}

func (c *customUserConversationModel) DeleteByConversationIdAndUserId(ctx context.Context, conversationId int64, userId int64) error {
	query := fmt.Sprintf("delete from %s where `user_id` = ? and `conversation_id` = ?", c.table)
	_, err := c.conn.ExecCtx(ctx, query, userId, conversationId)
	return err
}
