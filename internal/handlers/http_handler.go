package handlers

import (
	"github.com/GeorgeShibanin/hezzl_test5/internal/storage"
	"time"
)

type HTTPHandler struct {
	storage storage.Storage
}

func NewHTTPHandler(storage storage.Storage) *HTTPHandler {
	return &HTTPHandler{
		storage: storage,
	}
}

type PostNameData struct {
	Name string `json:"name"`
}

type PatchData struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ResponseData struct {
	Id          int       `json:"id"`
	CampaignId  int       `json:"campaignId"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Priority    int       `json:"priority"`
	Removed     bool      `json:"removed"`
	CreatedAt   time.Time `json:"createdAt"`
}

type ResponseAllItems struct {
	Items []storage.Item
}

type ResponseDeleteItem struct {
	Id         int  `json:"id"`
	CampaignId int  `json:"campaignId"`
	Removed    bool `json:"removed"`
}
