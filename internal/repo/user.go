package repo

import (
	"github.com/jmoiron/sqlx"
)

type UserRepo struct {
	Db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) *UserRepo {
	return &UserRepo{
		Db: db,
	}
}
