package model

import (
	"context"
	"errors"
	"testing"
)

// Mock repo and interface method implementations to be used for testing

func NewMockRepo() Repository {
	return Repository{
		Campaign: &MockCampaignRepo{
			Campaigns: map[int]*Campaign{
				1: {
					Name:          "ford",
					StartDate:     "2025-10-12",
					EndDate:       "2026-01-01",
					TargetDmaId:   501,
					AdId:          2,
					AdName:        "ForBiggerEscapes",
					AdDuration:    15,
					AdCreativeId:  102,
					AdCreativeUrl: "http://commondatastorage.googleapis.com/gtv-videos-bucket/sample/ForBiggerEscapes.mp4",
				},
			},
		},
	}
}

type MockCampaignRepo struct {
	Campaigns map[int]*Campaign
}

func (m *MockCampaignRepo) addMockCampaigns(t *testing.T, campaigns ...*Campaign) {
	for _, c := range campaigns {
		id := 1
		m.Campaigns[id] = c
		id++
	}
}

func (m *MockCampaignRepo) Create(ctx context.Context, campaign *Campaign) error {
	return nil
}

func (m *MockCampaignRepo) Delete(ctx context.Context, campaignId int64) error {
	return nil
}

func (m *MockCampaignRepo) Update(ctx context.Context, campaignId int64, campaign *Campaign) error {
	return nil
}

func (m *MockCampaignRepo) GetAll(ctx context.Context) ([]Campaign, error) {
	return nil, nil
}

func (m *MockCampaignRepo) GetById(ctx context.Context, campaignId int64) (*Campaign, error) {
	if campaign, ok := m.Campaigns[int(campaignId)]; ok {
		return campaign, nil
	}
	return nil, errors.New("campaign not found in mock repo")
}

func (m *MockCampaignRepo) GetByDma(ctx context.Context, campaignId int64) (*Campaign, error) {
	return nil, nil
}
