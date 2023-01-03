package router

import (
	"net/http"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"tzh.com/web/handler/check"
	"tzh.com/web/handler/user"
	"tzh.com/web/router/middleware"

	// docs is generated by Swag CLI, you have to import it.
	_ "tzh.com/web/docs"
)

// Load 载入中间件
func Load(g *gin.Engine, mw ...gin.HandlerFunc) *gin.Engine {
	g.Use(gin.Logger())
	g.Use(gin.Recovery())
	g.Use(middleware.NoCache())
	g.Use(middleware.Options())
	g.Use(middleware.Secure())
	g.Use(mw...)

	// 404 handler
	g.NoRoute(func(ctx *gin.Context) {
		ctx.String(http.StatusNotFound, "incorrect api router")
	})

	// pprof router, default is "/debug/pprof"
	pprof.Register(g)

	// swagger 文档
	// The url pointing to API definition
	// /swagger/index.html
	url := ginSwagger.URL("/swagger/doc.json")
	g.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	// 健康检查
	g.GET("/health", check.HealthCheck)
	// 创建用户, 无验证版
	g.POST("/user", user.Create)

	// 版本 v1
	g.POST("/v1/login", user.Login)

	g.POST("/v1/create", user.Create) // 为了方便创建用户, 无需认证

	u := g.Group("/v1/user")
	u.Use(middleware.AuthJWT()) // 添加认证
	{
		u.GET("", user.List)
		u.POST("", user.Create)
		u.GET("/:id", user.Get)
		u.PUT("/:id", user.Save)
		u.PATCH("/:id", user.Update)
		u.DELETE("/:id", user.Delete)
	}

	um := g.Group("/v1/username")
	um.Use(middleware.AuthJWT())
	{
		um.GET("/:name", user.GetByName)
	}

	checkRoute := g.Group("/v1/check")
	{
		checkRoute.GET("/health", check.HealthCheck)
		checkRoute.GET("/disk", check.DiskCheck)
		checkRoute.GET("/cpu", check.CPUCheck)
		checkRoute.GET("/memory", check.MemoryCheck)
	}

	return g

}