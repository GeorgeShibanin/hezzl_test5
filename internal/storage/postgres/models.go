package postgres

import "time"

type Item struct {
	Id          int       `json:"id"`
	CampaignId  int       `json:"campaignId"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Priority    int       `json:"priority"`
	Removed     bool      `json:"removed"`
	CreatedAt   time.Time `json:"createdAt"`
}
