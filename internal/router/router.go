package router

import (
	"net/http"
	"post/internal/auth"
	"post/internal/dashboard"
	"post/internal/pkg/cache"
	"post/internal/pkg/config"
	"post/internal/pkg/database"
	"post/internal/pkg/middleware"
	"post/internal/pkg/response"
	"post/internal/post"
	"post/internal/profile"
	"post/internal/user"

	"github.com/gin-gonic/gin"
)

func Init(cfg *config.Config) *gin.Engine {
	// Initialize Gin
	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()

	// Global Middleware
	r.Use(middleware.RequestID())
	r.Use(middleware.Logger())
	r.Use(gin.Recovery()) // Using Gin's default recovery or custom one? Plan said custom.
	// Let's use our custom recovery if implemented, but I see I implemented `middleware.Recovery()`.
	// Overriding gin.Recovery with ours.
	// Actually, `gin.New()` doesn't have default middleware.
	r.Use(middleware.Recovery())

	// Dependencies
	db := database.GetDB()

	// Repositories
	userRepo := user.NewRepository(db)
	profileRepo := profile.NewRepository(db)
	postRepo := post.NewRepository(db)

	// Services
	jwtService := auth.NewJWTService(cfg)
	authService := auth.NewService(userRepo, jwtService)
	userService := user.NewService(userRepo)
	profileService := profile.NewService(profileRepo)

	// Initialize Cache (100 items)
	postCache, err := cache.NewLRUCache(100)
	if err != nil {
		panic(err)
	}
	postService := post.NewService(postRepo, postCache)
	// Wait, Check post service implementation. It only took repo. Good.

	// Handlers
	authHandler := auth.NewHandler(authService)
	userHandler := user.NewHandler(userService)
	profileHandler := profile.NewHandler(profileService)
	postHandler := post.NewHandler(postService)

	// Auth Middleware
	authMiddleware := auth.Middleware(jwtService)

	// Routes
	api := r.Group("/api")
	{
		// Auth
		authRoutes := api.Group("/auth")
		{
			authRoutes.POST("/signup", authHandler.Signup)
			authRoutes.POST("/signin", authHandler.Signin)
		}

		// User
		userRoutes := api.Group("/users")
		userRoutes.Use(authMiddleware)
		{
			userRoutes.GET("/me", userHandler.GetProfile)
			userRoutes.GET("/:id", userHandler.GetUserByID)
		}

		// Profile
		profileRoutes := api.Group("/profile")
		profileRoutes.Use(authMiddleware)
		{
			profileRoutes.GET("/", profileHandler.GetProfile)
			profileRoutes.PUT("/", profileHandler.UpsertProfile)
		}

		// Post
		postRoutes := api.Group("/posts")
		{
			postRoutes.GET("/", postHandler.GetAllPosts)
			postRoutes.GET("/:id", postHandler.GetPostByID)

			// Protected
			postRoutes.Use(authMiddleware)
			postRoutes.POST("/", postHandler.CreatePost)
		}
	}

	// Load Templates
	r.LoadHTMLGlob("web/templates/**/*")

	// Admin Dashboard
	dashboardHandler := dashboard.NewHandler(userService, postService)
	admin := r.Group("/admin")
	admin.Use(gin.BasicAuth(gin.Accounts{
		cfg.App.AdminUser: cfg.App.AdminPassword,
	}))
	{
		admin.GET("/", dashboardHandler.ServeIndex)
		admin.GET("/users", dashboardHandler.ServeUsers)
		admin.GET("/posts", dashboardHandler.ServePosts)

		// Actions
		admin.DELETE("/users/:id", dashboardHandler.DeleteUser)
		admin.DELETE("/posts/:id", dashboardHandler.DeletePost)
	}

	// 404 Handler
	r.NoRoute(func(c *gin.Context) {
		response.Error(c, http.StatusNotFound, "Route not found", nil)
	})

	return r
}
