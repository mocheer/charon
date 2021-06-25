package dmap

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"

	"github.com/mocheer/charon/mw"
)

// Use 初始化 dmap 路由
func Use(api fiber.Router) {
	router := api.Group("/dmap")
	router.Get("/image/:id/:z", mw.NewLimiter(limiter.Config{Max: 1, Expiration: 30 * time.Second}), createImageHandle)
	router.Get("/image/:id/:z/:y/:x", imageTileHandle)
	//
	router.Get("/layer/:id", mw.Cache, layerHandle)
	//
	router.Get("/feature/:id", mw.Cache, featureHandle)
	router.Get("/feature2/:id", mw.Cache, featureHandle2)
	//
	router.Get("/identify/:id", mw.Cache, identifyHandle)
	//
	router.Get("/kriging/grid", krigingGridHandle)
	router.Post("/kriging/grid", krigingGridHandle)
	//
	router.Get("/kriging/image", krigingImageHandle)
	router.Post("/kriging/image", krigingImageHandle)
}
