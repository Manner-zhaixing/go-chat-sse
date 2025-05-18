package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ ConversationModel = (*customConversationModel)(nil)

type (
	// ConversationModel is an interface to be customized, add more methods here,
	// and implement the added methods in customConversationModel.
	ConversationModel interface {
		conversationModel
		withSession(session sqlx.Session) ConversationModel
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

func (m *customConversationModel) withSession(session sqlx.Session) ConversationModel {
	return NewConversationModel(sqlx.NewSqlConnFromSession(session))
}
