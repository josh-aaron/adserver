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

func (r *CampaignRepo) GetAll(ctx context.Context) ([]Campaign, error) {
	log.Println("campaign.GetAll()")
	query := `SELECT * FROM campaign`
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	rows, err := r.db.QueryContext(ctx, query)
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

	log.Printf("campaign.GetAll() returning all campaigns: %v", campaigns)

	return campaigns, nil
}

func (r *CampaignRepo) Delete(ctx context.Context, campaignId int64) error {
	log.Println("campaign.Delete()")
	query := `DELETE FROM campaign WHERE id = $1`
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	result, err := r.db.ExecContext(ctx, query, campaignId)
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
	log.Printf("campaign.Delete: successfully deleted campaign id %d", campaignId)
	return nil
}

func (r *CampaignRepo) Create(ctx context.Context, campaign *Campaign) error {
	log.Println("campaign.Create()")
	query := `
	INSERT INTO campaign (name, start_date, end_date, target_dma_id, ad_id, ad_name, ad_duration, ad_creative_id, ad_creative_url)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) 
	RETURNING id
	`

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	err := r.db.QueryRowContext(
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

	log.Printf("campaign.Create() sucessfully created campaign: %v", campaign)

	return nil
}

func (r *CampaignRepo) Update(ctx context.Context, campaignId int64, campaign *Campaign) error {
	log.Println("campaign.Update()")
	query := `
		UPDATE campaign
		SET name = $1, start_date = $2, end_date = $3, target_dma_id = $4, ad_id = $5, ad_name = $6, ad_duration = $7, ad_creative_id = $8, ad_creative_url = $9
		WHERE id = $10
		RETURNING id
	`
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	err := r.db.QueryRowContext(
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
		campaignId,
	).Scan(&campaign.Id)
	if err != nil {
		log.Println(err)
		switch err {
		case sql.ErrNoRows:
			return ErrNotFound
		default:
			return err
		}
	}

	log.Printf("campaign.Create() sucessfully updated campaign: %v", campaign)

	return nil

}

func (r *CampaignRepo) GetById(ctx context.Context, campaignId int64) (*Campaign, error) {
	log.Println("campaign.GetById()")

	query := `SELECT * FROM campaign WHERE id = $1`
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	var campaign Campaign
	err := r.db.QueryRowContext(ctx, query, campaignId).Scan(
		&campaign.Id,
		&campaign.Name,
		&campaign.StartDate,
		&campaign.EndDate,
		&campaign.TargetDmaId,
		&campaign.AdId,
		&campaign.AdName,
		&campaign.AdDuration,
		&campaign.AdCreativeId,
		&campaign.AdCreativeUrl,
	)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	log.Printf("campaign.GetById() returning campaign: %v", campaign)

	return &campaign, nil
}

func (r *CampaignRepo) GetByDma(ctx context.Context, campaignId int64) (*Campaign, error) {
	log.Println("campaign.GetByDma()")
	query := `SELECT * FROM campaign WHERE target_dma_id = $1`
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	var campaign Campaign
	err := r.db.QueryRowContext(ctx, query, campaignId).Scan(
		&campaign.Id,
		&campaign.Name,
		&campaign.StartDate,
		&campaign.EndDate,
		&campaign.TargetDmaId,
		&campaign.AdId,
		&campaign.AdName,
		&campaign.AdDuration,
		&campaign.AdCreativeId,
		&campaign.AdCreativeUrl,
	)

	if err != nil {
		return nil, err
	}

	log.Printf("campaign.GetByDma() returning campaign: %v", campaign)

	return &campaign, nil
}
