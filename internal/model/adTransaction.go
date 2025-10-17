package model

import (
	"context"
	"database/sql"
	"log"
	"time"
)

// Currently unused - in the future we may need to modify an AdTransaction based on business logic
type AdTransaction struct {
	TransactionId int64
	VastResponse  string
	AdResponse    []byte
	ClientDmaId   int64
	CampaignId    int64
}

type AdTransactionRepo struct {
	db *sql.DB
}

func (r *AdTransactionRepo) CreateTransactionId() int64 {
	return time.Now().UnixMilli()
}

func (r *AdTransactionRepo) CreateAdTransaction(ctx context.Context, transactionId int64, adrequest string, vastXml []byte, clientDmaId int64, campaignId int64) {
	log.Printf("adTransaction.CreateAdTransaction() for transactionID: %v", transactionId)

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
