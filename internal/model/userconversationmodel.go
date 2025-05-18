package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ UserConversationModel = (*customUserConversationModel)(nil)

type (
	// UserConversationModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserConversationModel.
	UserConversationModel interface {
		userConversationModel
		withSession(session sqlx.Session) UserConversationModel
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

func (m *customUserConversationModel) withSession(session sqlx.Session) UserConversationModel {
	return NewUserConversationModel(sqlx.NewSqlConnFromSession(session))
}
