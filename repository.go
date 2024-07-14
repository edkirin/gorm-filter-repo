package repository

import (
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

const DEFAULT_ID_FIELD = "id"

type MethodInitInterface[T schema.Tabler] interface {
	Init(repo *RepoBase[T])
}

type RepoOptions struct {
	IdField string
}

type RepoBase[T schema.Tabler] struct {
	IdField string
	dbConn  *gorm.DB

	ListMethod[T]
	GetMethod[T]
	ExistsMethod[T]
	CountMethod[T]
	SaveMethod[T]
	UpdateMethod[T]
	DeleteMethod[T]
}

func (repo *RepoBase[T]) InitMethods(methods []MethodInitInterface[T]) {
	for _, method := range methods {
		method.Init(repo)
	}
}

func (m *RepoBase[T]) Init(dbConn *gorm.DB, options *RepoOptions) {
	m.dbConn = dbConn

	if options == nil {
		// no options provided? set defaults
		m.IdField = DEFAULT_ID_FIELD
	} else {
		if len(options.IdField) > 0 {
			m.IdField = options.IdField
		}
	}

	methods := []MethodInitInterface[T]{
		&m.ListMethod,
		&m.GetMethod,
		&m.ExistsMethod,
		&m.CountMethod,
		&m.SaveMethod,
		&m.UpdateMethod,
		&m.DeleteMethod,
	}
	m.InitMethods(methods)
}
