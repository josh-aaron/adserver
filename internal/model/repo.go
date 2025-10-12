package model

import (
	"context"
	"database/sql"
)

type Repository struct {
	Campaign interface {
		Create(context.Context, *Campaign) error
	}
}

func NewRepository(db *sql.DB) Repository {
	return Repository{
		Campaign: &CampaignRepo{db},
	}

}
