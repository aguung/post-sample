package dashboard

import (
	"net/http"
	"strconv"

	"post/internal/pkg/response"
	"post/internal/post"
	"post/internal/user"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	userService user.Service
	postService post.Service
}

func NewHandler(userService user.Service, postService post.Service) *Handler {
	return &Handler{userService, postService}
}

func (h *Handler) ServeIndex(c *gin.Context) {
	users, _ := h.userService.GetAll()
	posts, _ := h.postService.GetAll()

	data := gin.H{
		"Page":       "overview",
		"TotalUsers": len(users),
		"TotalPosts": len(posts),
		"CurrentURL": "/admin/",
	}
	c.HTML(http.StatusOK, "base.html", data)
}

func (h *Handler) ServeUsers(c *gin.Context) {
	users, err := h.userService.GetAll()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "users.html", gin.H{"Error": err.Error()})
		return
	}

	data := gin.H{
		"Page":       "users",
		"Users":      users,
		"CurrentURL": "/admin/users",
	}
	c.HTML(http.StatusOK, "base.html", data)
}

func (h *Handler) ServePosts(c *gin.Context) {
	posts, err := h.postService.GetAll()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "posts.html", gin.H{"Error": err.Error()})
		return
	}

	data := gin.H{
		"Page":       "posts",
		"Posts":      posts,
		"CurrentURL": "/admin/posts",
	}
	c.HTML(http.StatusOK, "base.html", data)
}

func (h *Handler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid ID", err.Error())
		return
	}

	if err := h.userService.Delete(uint(id)); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to delete user", err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) DeletePost(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid ID", err.Error())
		return
	}

	if err := h.postService.Delete(uint(id)); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to delete post", err.Error())
		return
	}

	c.Status(http.StatusOK)
}
