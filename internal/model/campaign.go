package model

import (
	"context"
	"database/sql"
	"log"
)

type Campaign struct {
	Id            int64  `json:"id"`
	Name          string `json:"name"`
	StartDate     string `json:"startDate"`
	EndDate       string `json:"endDate"`
	TargetDmaId   int    `json:"targetDmaId"`
	AdId          int    `json:"adId"`
	AdName        string `json:"adName"`
	AdDuration    int    `json:"adDuration"`
	AdCreativeId  int    `json:"adCreativeId"`
	AdCreativeUrl string `json:"adCreativeUrl"`
}

type CampaignRepo struct {
	db *sql.DB
}

func (s *CampaignRepo) Create(ctx context.Context, campaign *Campaign) error {
	log.Println("campaign.Create()")
	query := `
	INSERT INTO campaign (name, start_date, end_date, target_dma_id, ad_id, ad_name, ad_duration, ad_creative_id, ad_creative_url)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id
	`
	err := s.db.QueryRowContext(
		ctx,
		query,
		campaign.Name,
		campaign.StartDate,
		campaign.EndDate,
		campaign.TargetDmaId,
		campaign.AdId,
		campaign.AdName,
		campaign.AdDuration,
		campaign.AdCreativeId,
		campaign.AdCreativeUrl,
	).Scan(
		&campaign.Id,
	)
	if err != nil {
		return err
	}

	return nil
}
