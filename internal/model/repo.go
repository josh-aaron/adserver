package model

import (
	"context"
	"database/sql"
)

type Repository struct {
	Campaign interface {
		Create(context.Context, *Campaign) error
		Delete(context.Context, int64) error
	}
}

func NewRepository(db *sql.DB) Repository {
	return Repository{
		Campaign: &CampaignRepo{db},
	}

}
