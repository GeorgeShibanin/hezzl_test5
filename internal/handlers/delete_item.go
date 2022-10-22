package handlers

import (
	"encoding/json"
	"github.com/GeorgeShibanin/hezzl_test5/internal/storage"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
)

func (h *HTTPHandler) HandleDeleteItem(rw http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(rw, "invalid query params", http.StatusBadRequest)
		return
	}
	newId, err := strconv.Atoi(id)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	campaignId, err := strconv.Atoi(r.URL.Query().Get("campaignId"))
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	item, err := h.storage.DeleteItem(r.Context(), storage.Id(newId), storage.CampaignId(campaignId))
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	response := ResponseDeleteItem{
		Id:         item.Id,
		CampaignId: item.CampaignId,
		Removed:    item.Removed,
	}
	rawResponse, err := json.Marshal(response)
	if err != nil {
		err = errors.Wrap(err, "can't marshall response")
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = rw.Write(rawResponse)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
}
