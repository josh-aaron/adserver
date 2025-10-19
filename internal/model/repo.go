package model

import (
	"context"
	"database/sql"
)

type Repository struct {
	Campaign interface {
		Create(context.Context, *Campaign) error
		Delete(context.Context, int64) error
		Update(context.Context, int64, *Campaign) error
		GetAll(context.Context) ([]Campaign, error)
		GetById(context.Context, int64) (*Campaign, error)
		GetByDma(context.Context, int64) (*Campaign, error)
	}
	VastResponse interface {
		GetVast(context.Context, *Campaign, int, int64) (*VAST, int, error)
	}
	AdTransaction interface {
		CreateTransactionId() int64
		LogAdTransaction(context.Context, int64, string, []byte, int64, int64)
		LogBeacons(context.Context, int64, string, string) error
		GetAllAdTransactions(context.Context) ([]AdTransaction, error)
		GetBeaconsByTransactionId(context.Context, int64) ([]AdBeacon, error)
	}
}

func NewRepository(db *sql.DB) Repository {
	return Repository{
		Campaign:      &CampaignRepo{db},
		VastResponse:  &VastResponseRepo{db},
		AdTransaction: &AdTransactionRepo{db},
	}

}
