package model

import (
	"context"
	"database/sql"
)

type Campaign struct {
	Id            int64  `json:"id"`
	Name          string `json:"name"`
	StartDate     string `json:"startDate"`
	EndDate       string `json:"endDate"`
	TargetDmaId   int64  `json:"targetDmaId"`
	AdId          int64  `json:"adId"`
	AdName        string `json:"adName"`
	AdDuration    int64  `json:"adDuration"`
	AdCreativeId  int64  `json:"adCreativeId"`
	AdCreativeUrl string `json:"adCreativeUrl"`
}

type CampaignRepo struct {
	db *sql.DB
}

func (s *CampaignRepo) Create(ctx context.Context, campaign *Campaign) error {
	query := `
	INSERT INTO campaign (name, startDate, endDate, targetDmaId, adId, adName, adDuration, adCreativeId, adCreativeUrl)
	VALUES ($1, $2, $3, $4) RETURNING id
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
