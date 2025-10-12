package model

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"
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

func (s *CampaignRepo) GetAll(ctx context.Context) ([]Campaign, error) {
	log.Println("campaign.GetAll()")
	query := `SELECT * FROM campaign`
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var campaigns []Campaign
	for rows.Next() {
		var c Campaign
		err := rows.Scan(
			&c.Id,
			&c.Name,
			&c.StartDate,
			&c.EndDate,
			&c.TargetDmaId,
			&c.AdId,
			&c.AdName,
			&c.AdDuration,
			&c.AdCreativeId,
			&c.AdCreativeUrl,
		)
		if err != nil {
			return nil, err
		}
		campaigns = append(campaigns, c)
	}

	return campaigns, nil
}

func (s *CampaignRepo) Delete(ctx context.Context, campaignId int64) error {
	log.Println("campaign.Delete()")
	query := `DELETE FROM campaign WHERE id = $1`
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	result, err := s.db.ExecContext(ctx, query, campaignId)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return nil
	}
	if rows == 0 {
		return errors.New("campaign not found")
	}
	log.Printf("campaign.Delete: successfully delete campaign with id %d", campaignId)
	return nil
}

func (s *CampaignRepo) Create(ctx context.Context, campaign *Campaign) error {
	log.Println("campaign.Create()")
	query := `
	INSERT INTO campaign (name, start_date, end_date, target_dma_id, ad_id, ad_name, ad_duration, ad_creative_id, ad_creative_url)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id
	`

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

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
