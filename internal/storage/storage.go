package storage

import (
	"context"
	"errors"
	"fmt"
	"time"
)

var (
	StorageError    = errors.New("storage")
	ErrCollision    = fmt.Errorf("%w.collision", StorageError)
	ErrAlreadyExist = errors.New("key already exist")
	ErrNotFound     = fmt.Errorf("%w.not_found", StorageError)
)

type Id int
type Name string
type CampaignId int
type Description string

type Item struct {
	Id          int       `json:"id"`
	CampaignId  int       `json:"campaign_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Priority    int       `json:"priority"`
	Removed     bool      `json:"removed"`
	CreatedAt   time.Time `json:"created_at"`
}

type Storage interface {
	PostItem(ctx context.Context, campaignId CampaignId, name Name) (Item, error)
	PatchItem(ctx context.Context, id Id, campaignId CampaignId, name Name, description Description) (Item, error)
	GetItems(ctx context.Context) ([]Item, error)
	DeleteItem(ctx context.Context, id Id, campaignId CampaignId) (Item, error)
}
