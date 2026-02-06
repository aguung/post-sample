package post

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

func (h *Handler) CreatePost(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	var input CreatePostInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid input", err)
		return
	}

	post, err := h.service.Create(userID, input)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to create post", err.Error())
		return
	}

	response.Success(c, http.StatusCreated, "Post created successfully", post)
}

func (h *Handler) GetAllPosts(c *gin.Context) {
	posts, err := h.service.GetAll()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to fetch posts", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Posts retrieved", posts)
}

func (h *Handler) GetPostByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid ID", err.Error())
		return
	}

	post, err := h.service.GetByID(uint(id))
	if err != nil {
		response.Error(c, http.StatusNotFound, "Post not found", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Post retrieved", post)
}
