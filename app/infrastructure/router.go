package infrastructure

import (
	"database/sql"
	"net/http"
	"os"

	"github.com/bmf-san/gobel-api/app/interfaces"
	"github.com/bmf-san/gobel-api/app/middleware"
	"github.com/bmf-san/gobel-api/app/usecases"
	"github.com/bmf-san/goblin"
	"github.com/go-redis/redis/v7"
)

// Dispatch handle routing
func Dispatch(connMySQL *sql.DB, connRedis *redis.Client, logger usecases.Logger) {
	asset := middleware.NewAsset(connRedis, logger)

	publicMws := middleware.NewMiddlewares(asset.AccessLog, asset.CORS)
	privateMws := middleware.NewMiddlewares(asset.AccessLog, asset.CORS, asset.Auth)

	authController := interfaces.NewAuthController(connMySQL, connRedis, logger)
	postController := interfaces.NewPostController(connMySQL, connRedis, logger)
	commentController := interfaces.NewCommentController(connMySQL, logger)
	categoryController := interfaces.NewCategoryController(connMySQL, logger)
	tagController := interfaces.NewTagController(connMySQL, logger)

	r := goblin.NewRouter()

	r.GET("/posts", publicMws.Then(postController.Index))
	r.GET("/posts/categories/:name", publicMws.Then(postController.IndexByCategory))
	r.GET("/posts/tags/:name", publicMws.Then(postController.IndexByTag))
	r.GET("/posts/:title", publicMws.Then(postController.Show))
	r.POST("/posts/:title/comments", publicMws.Then(commentController.Store))

	r.GET("/categories", publicMws.Then(categoryController.Index))
	r.GET("/categories/:name", publicMws.Then(categoryController.Show))

	r.GET("/tags", publicMws.Then(tagController.Index))
	r.GET("/tags/:name", publicMws.Then(tagController.Show))

	r.POST("/signin", publicMws.Then(authController.SignIn))

	r.POST("/private/signin", privateMws.Then(authController.SignOut))

	r.GET("/private/posts", privateMws.Then(postController.IndexPrivate))
	r.GET("/private/posts/:id", privateMws.Then(postController.ShowPrivate))
	r.POST("/private/posts", privateMws.Then(postController.StorePrivate))
	r.PATCH("/private/posts/:id", privateMws.Then(postController.UpdatePrivate))
	r.DELETE("/private/posts/:id", privateMws.Then(postController.DestroyPrivate))

	r.GET("/private/comments", privateMws.Then(commentController.IndexPrivate))
	r.GET("/private/comments/:id", privateMws.Then(commentController.ShowPrivate))
	r.PATCH("/private/comments/:id/status", privateMws.Then(commentController.UpdateStatusPrivate))

	r.GET("/private/categories", privateMws.Then(categoryController.IndexPrivate))
	r.GET("/private/categories/:id", privateMws.Then(categoryController.ShowPrivate))
	r.POST("/private/categories", privateMws.Then(categoryController.StorePrivate))
	r.PATCH("/private/categories/:id", privateMws.Then(categoryController.UpdatePrivate))
	r.DELETE("/private/categories/:id", privateMws.Then(categoryController.DestroyPrivate))

	r.GET("/private/tags", privateMws.Then(tagController.IndexPrivate))
	r.GET("/private/tags/:id", privateMws.Then(tagController.ShowPrivate))
	r.POST("/private/tags", privateMws.Then(tagController.StorePrivate))
	r.PATCH("/private/tags/:id", privateMws.Then(tagController.UpdatePrivate))
	r.DELETE("/private/tags/:id", privateMws.Then(tagController.DestroyPrivate))

	if err := http.ListenAndServe(":"+os.Getenv("SERVER_PORT"), r); err != nil {
		logger.LogError(err)
	}
}
