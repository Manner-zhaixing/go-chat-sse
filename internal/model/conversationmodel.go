package model

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ ConversationModel = (*customConversationModel)(nil)

type (
	// ConversationModel is an interface to be customized, add more methods here,
	// and implement the added methods in customConversationModel.
	ConversationModel interface {
		conversationModel
		FindOneBySessionId(ctx context.Context, id int64) (*Conversation, error)
	}

	customConversationModel struct {
		*defaultConversationModel
	}
)

// NewConversationModel returns a model for the database table.
func NewConversationModel(conn sqlx.SqlConn) ConversationModel {
	return &customConversationModel{
		defaultConversationModel: newConversationModel(conn),
	}
}

func (m *customConversationModel) FindOneBySessionId(ctx context.Context, sessionid int64) (*Conversation, error) {
	var resp Conversation
	query := fmt.Sprintf("select %s from %s where `session_id` = ? limit 1", messageRows, m.table)
	err := m.conn.QueryRowCtx(ctx, &resp, query, sessionid)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}
