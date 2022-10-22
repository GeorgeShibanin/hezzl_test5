package handlers

import (
	"encoding/json"
	"github.com/GeorgeShibanin/hezzl_test5/internal/storage"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
)

func (h *HTTPHandler) HandlePostItem(rw http.ResponseWriter, r *http.Request) {
	campaignId, err := strconv.Atoi(r.URL.Query().Get("campaignId"))
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	var name PostNameData
	err = json.NewDecoder(r.Body).Decode(&name)

	item, err := h.storage.PostItem(r.Context(), storage.CampaignId(campaignId), storage.Name(name.Name))
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	response := ResponseData(item)
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
