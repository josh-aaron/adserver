package model

import (
	"context"
	"database/sql"
	"log"
	"time"
)

type AdTransaction struct {
	TransactionId int64
	AdRequest     string
	VastResponse  string
	ClientDmaId   int64
	CampaignId    int64
}

type AdBeacon struct {
	Id            int64
	TransactionId int64
	BeaconUrl     string
	BeaconName    string
}

type AdTransactionRepo struct {
	db *sql.DB
}

func (r *AdTransactionRepo) CreateTransactionId() int64 {
	return time.Now().UnixMilli()
}

func (r *AdTransactionRepo) LogAdTransaction(ctx context.Context, transactionId int64, adrequest string, vastXml []byte, clientDmaId int64, campaignId int64) {
	log.Printf("adTransaction.LogAdTransaction() for transactionID: %v", transactionId)

	query := `
	INSERT INTO ad_transaction (transaction_id, ad_request, vast_response, client_dma_id, campaign_id)
	VALUES ($1, $2, $3, $4, $5)
	`

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	_, err := r.db.ExecContext(
		ctx,
		query,
		transactionId,
		adrequest,
		vastXml,
		clientDmaId,
		campaignId,
	)
	if err != nil {
		log.Println(err)
		return
	}
}

func (r *AdTransactionRepo) LogBeacon(ctx context.Context, transactionId int64, beaconUri string, beaconName string) error {
	log.Printf("adTransaction.LogBeacon() logging %v for transactionID: %v", beaconName, transactionId)
	query := `
	INSERT INTO ad_beacon (transaction_id, beacon_url, beacon_name)
	VALUES ($1, $2, $3)
	`

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	_, err := r.db.ExecContext(
		ctx,
		query,
		transactionId,
		beaconUri,
		beaconName,
	)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (r *AdTransactionRepo) GetAllAdTransactions(ctx context.Context) ([]AdTransaction, error) {
	log.Println("AdTransactionRepo.GetAllAdTransactions()")

	query := `SELECT transaction_id, ad_request, vast_response, client_dma_id, campaign_id FROM ad_transaction`
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var adTransactions []AdTransaction
	for rows.Next() {
		var at AdTransaction
		err := rows.Scan(
			&at.TransactionId,
			&at.AdRequest,
			&at.VastResponse,
			&at.ClientDmaId,
			&at.CampaignId,
		)
		if err != nil {
			return nil, err
		}
		adTransactions = append(adTransactions, at)
	}

	log.Printf("AdTransactionRepo.GetAllAdTransactions() returning %v adTransactions", len(adTransactions))
	return adTransactions, nil
}

func (r *AdTransactionRepo) GetBeaconsByTransactionId(ctx context.Context, transactionId int64) ([]AdBeacon, error) {
	log.Printf("AdTransactionRepo.GetBeaconsByTransactionId() for transactionId %v", transactionId)

	query := `
	SELECT id, transaction_id, beacon_url, beacon_name 
	FROM ad_beacon 
	WHERE transaction_id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	rows, err := r.db.QueryContext(ctx, query, transactionId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var beacons []AdBeacon
	for rows.Next() {
		var ab AdBeacon
		err := rows.Scan(
			&ab.Id,
			&ab.TransactionId,
			&ab.BeaconUrl,
			&ab.BeaconName,
		)
		if err != nil {
			return nil, err
		}
		beacons = append(beacons, ab)
	}
	if len(beacons) == 0 {
		return nil, ErrNotFound
	}

	log.Printf("AdTransactionRepo.GetBeaconsByTransactionId() returning beacons: %v", beacons)
	return beacons, nil
}
