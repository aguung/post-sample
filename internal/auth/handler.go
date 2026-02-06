package auth

import (
	"errors"
	"net/http"

	pkgdb "post/internal/pkg/database"
	"post/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service}
}

func (h *Handler) Signup(c *gin.Context) {
	var input SignupInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid input", err)
		return
	}

	user, err := h.service.Signup(input)
	if err != nil {
		if errors.Is(err, pkgdb.ErrDuplicateKey) {
			response.Error(c, http.StatusUnprocessableEntity, "Email already registered", nil)
			return
		}
		response.Error(c, http.StatusInternalServerError, "Failed to create user", err.Error())
		return
	}

	response.Success(c, http.StatusCreated, "User registered successfully", user)
}

func (h *Handler) Signin(c *gin.Context) {
	var input SigninInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid input", err)
		return
	}

	token, refreshToken, err := h.service.Signin(input)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "Login failed", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Login successful", gin.H{
		"token":         token,
		"refresh_token": refreshToken,
	})
}
