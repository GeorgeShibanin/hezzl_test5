package handlers

import (
	"encoding/json"
	"github.com/pkg/errors"
	"net/http"
)

func (h *HTTPHandler) HandleGetItem(rw http.ResponseWriter, r *http.Request) {

	items, err := h.storage.GetItems(r.Context())
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	var response ResponseAllItems
	response.Items = items
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
