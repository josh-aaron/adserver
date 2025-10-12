package model

import (
	"context"
	"database/sql"
)

type Repository struct {
	Campaign interface {
		Create(context.Context, *Campaign) error
		Delete(context.Context, int64) error
		Update(context.Context, int64) error
		GetAll(context.Context) ([]Campaign, error)
		GetById(context.Context, int64) (*Campaign, error)
	}
}

func NewRepository(db *sql.DB) Repository {
	return Repository{
		Campaign: &CampaignRepo{db},
	}

}
