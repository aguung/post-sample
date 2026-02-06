package user

import (
	"net/http"
	"strconv"

	"post/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service}
}

func (h *Handler) GetProfile(c *gin.Context) {
	// userID is guaranteed by Auth middleware
	id := c.MustGet("userID").(uint)

	user, err := h.service.GetByID(id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "User not found", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "User profile retrieved", user)
}

func (h *Handler) GetUserByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid ID", err.Error())
		return
	}

	user, err := h.service.GetByID(uint(id))
	if err != nil {
		response.Error(c, http.StatusNotFound, "User not found", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "User retrieved", user)
}
