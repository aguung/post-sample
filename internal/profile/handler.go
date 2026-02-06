package profile

import (
	"net/http"

	"post/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service}
}

func (h *Handler) UpsertProfile(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	var input ProfileInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid input", err)
		return
	}

	profile, err := h.service.CreateOrUpdate(userID, input)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to update profile", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Profile updated successfully", profile)
}

func (h *Handler) GetProfile(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	profile, err := h.service.GetByUserID(userID)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Profile not found", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Profile retrieved", profile)
}
