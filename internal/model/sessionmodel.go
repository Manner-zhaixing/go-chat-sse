package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ SessionModel = (*customSessionModel)(nil)

type (
	// SessionModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSessionModel.
	SessionModel interface {
		sessionModel
		withSession(session sqlx.Session) SessionModel
	}

	customSessionModel struct {
		*defaultSessionModel
	}
)

// NewSessionModel returns a model for the database table.
func NewSessionModel(conn sqlx.SqlConn) SessionModel {
	return &customSessionModel{
		defaultSessionModel: newSessionModel(conn),
	}
}

func (m *customSessionModel) withSession(session sqlx.Session) SessionModel {
	return NewSessionModel(sqlx.NewSqlConnFromSession(session))
}
