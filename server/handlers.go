package server

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/BeDCAF/vpn_service/panel"
)

type HTTPHandlers struct {
	client panel.VPNService
}

func NewHTTPHandlers(client panel.VPNService) *HTTPHandlers {
	return &HTTPHandlers{
		client: client,
	}
}

/*
pattern: /clients
method:  POST
info:    JSON in HTTP request body

succeed:
- status code:   201 Created
- response body: client config
failed:
- status code:   400, 502
*/
func (h *HTTPHandlers) HandleAddUser(w http.ResponseWriter, r *http.Request) {
	var user panel.Client

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		errDTO := ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}

		http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		return
	}

	if err := user.ValidateClient(); err != nil {
		errDTO := ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}

		http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	data, err := h.client.AddUser(ctx, panel.Response{Client: user, InboundIds: []int64{1}})
	if err != nil {
		errDTO := ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}

		http.Error(w, errDTO.ToString(), http.StatusBadGateway)
		return
	}

	body, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	_, _ = w.Write(body)
}

/*
pattern: /clients
method:  DELETE
info:    JSON in HTTP request body

succeed:
- status code:   204 No Content
failed:
- status code:   400, 502
*/
func (h *HTTPHandlers) HandleDeleteUser(w http.ResponseWriter, r *http.Request) {
	var user panel.Client

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		errDTO := ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}

		http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		return
	}

	if err := user.ValidateClient(); err != nil {
		errDTO := ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}

		http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	err := h.client.DeleteUser(ctx, panel.Response{Client: user, InboundIds: []int64{1}})
	if err != nil {
		errDTO := ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}

		http.Error(w, errDTO.ToString(), http.StatusBadGateway)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
