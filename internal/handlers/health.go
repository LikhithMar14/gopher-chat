package handlers

import (
	"net/http"

	"github.com/LikhithMar14/gopher-chat/internal/config"
	"github.com/LikhithMar14/gopher-chat/internal/utils"
)

const version = "1.0.0"

type HealthHandler struct {
	config config.Config
}

func NewHealthHandler(cfg config.Config) *HealthHandler {
	return &HealthHandler{
		config: cfg,
	}
}

func (h *HealthHandler) Handle(w http.ResponseWriter, r *http.Request) {
	data := utils.Envelope{
		"status":      "ok",
		"environment": h.config.Env,
		"version":     version,
	}

	if err := utils.WriteJSON(w, http.StatusOK, data); err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "Failed to write response")
	}
}
