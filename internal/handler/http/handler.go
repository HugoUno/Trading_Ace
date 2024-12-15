package http

import (
	"encoding/json"
	"net/http"

	"trading-ace/internal/service"

	"github.com/gorilla/mux"
)

type Handler struct {
	campaignService *service.CampaignService
}

func NewHandler(campaignService *service.CampaignService) *Handler {
	return &Handler{
		campaignService: campaignService,
	}
}

func (h *Handler) Router() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/api/v1/users/{address}/tasks", h.GetUserTasks).Methods("GET")
	r.HandleFunc("/api/v1/users/{address}/points/history", h.GetUserPointsHistory).Methods("GET")
	r.HandleFunc("/api/v1/leaderboard", h.GetLeaderboard).Methods("GET")

	return r
}

func (h *Handler) GetUserTasks(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["address"]

	tasks, err := h.campaignService.GetUserTaskStatus(r.Context(), address)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(tasks)
}

func (h *Handler) GetUserPointsHistory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["address"]

	tasks, err := h.campaignService.GetUserTaskStatus(r.Context(), address)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(tasks)
}

func (h *Handler) GetLeaderboard(w http.ResponseWriter, r *http.Request) {
	rankings, err := h.campaignService.GetLeaderboard(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rankings)
}
