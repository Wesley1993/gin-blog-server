package router

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"wuzhispace.com/config"
	"wuzhispace.com/internal/handler"
	"wuzhispace.com/internal/repository"
	"wuzhispace.com/middleware"
	pkgredis "wuzhispace.com/pkg/redis"
)

// writePermTable 写操作路由所需按钮权限标识声明表。
// 键格式：METHOD + 空格 + FullPath（路由注册模式）。未收录的路由不设 perm，
// RBACAuth 中间件对空 perm 默认放行（仅需登录），如登出、只读接口等对全体登录用户开放的操作。
// 站点配置/站点资料/常用网站/个人信息/文件上传/索引重建等写操作均已登记权限标识，
// 非超管用户须由角色 button_perms 显式授予，否则一律拒绝（如只读「测试」角色）。
var writePermTable = map[string]string{
	// 文章管理
	"POST /api/article/create":     "article:add",
	"PUT /api/article/update":      "article:edit",
	"DELETE /api/article/:id":      "article:delete",
	"POST /api/article/es/rebuild": "article:rebuild",
	// 分类管理
	"POST /api/category/create": "category:add",
	"PUT /api/category/update":  "category:edit",
	"DELETE /api/category/:id":  "category:delete",
	// 用户管理
	"POST /api/user/create":      "user:add",
	"PUT /api/user/update":       "user:edit",
	"DELETE /api/user/:id":       "user:delete",
	"PUT /api/user/resetPwd/:id": "user:resetPwd",
	"PUT /api/user/status/:id":   "user:status",
	// 角色管理
	"POST /api/role/create": "role:add",
	"PUT /api/role/update":  "role:edit",
	"DELETE /api/role/:id":  "role:delete",
	// 菜单管理
	"POST /api/menu/create": "menu:add",
	"PUT /api/menu/update":  "menu:edit",
	"DELETE /api/menu/:id":  "menu:delete",
	// 站点设置（基础配置保存 + OSS 保存/连通测试同属站点配置写操作）
	"PUT /api/site/config":   "site:edit",
	"POST /api/site/testOss": "site:edit",
	// 站长个人资料（站点资料页）
	"PUT /api/site/profile": "siteProfile:edit",
	// 常用网站（友情链接）
	"POST /api/site/links":       "link:add",
	"PUT /api/site/links/:id":    "link:edit",
	"DELETE /api/site/links/:id": "link:delete",
	// 个人中心：修改当前登录用户自己的昵称/头像/简介
	"PUT /api/profile/update": "profile:edit",
	// 文件上传（文章封面、分类图片、头像等上传统一入口）
	"POST /api/upload/oss": "upload:create",
}

// markPerm 权限标识前置中间件：按 METHOD+FullPath 从声明表查得所需权限写入 Context，
// 供其后的 middleware.RBACAuth 读取（因此必须挂载在 RBACAuth 之前）
func markPerm() gin.HandlerFunc {
	return func(c *gin.Context) {
		if perm, ok := writePermTable[c.Request.Method+" "+c.FullPath()]; ok {
			c.Set("perm", perm)
		}
		c.Next()
	}
}

// SetupRouter 注册所有 API 路由
func SetupRouter(
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
	roleHandler *handler.RoleHandler,
	menuHandler *handler.MenuHandler,
	categoryHandler *handler.CategoryHandler,
	articleHandler *handler.ArticleHandler,
	siteHandler *handler.SiteHandler,
	siteProfileHandler *handler.SiteProfileHandler,
	siteLinkHandler *handler.SiteLinkHandler,
	uploadHandler *handler.UploadHandler,
	profileHandler *handler.ProfileHandler,
	systemHandler *handler.SystemHandler,
	roleRepo *repository.RoleRepository,
	cfg *config.Config,
) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// CORS 跨域中间件（允许前端开发地址）
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:5174"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	redisClient := pkgredis.GetClient()

	// 健康检查接口（无需认证）
	router.GET("/health", func(c *gin.Context) {
		health := map[string]interface{}{
			"status": "ok",
			"time":   time.Now().Format("2006-01-02 15:04:05"),
		}
		c.JSON(http.StatusOK, health)
	})

	api := router.Group("/api")
	{
		// ==================== 认证接口（无需登录） ====================
		auth := api.Group("/auth")
		{
			auth.POST("/login", middleware.RateLimit(redisClient), authHandler.Login)
		}

		// ==================== 公开接口（博客前台使用，无需登录） ====================
		public := api.Group("/public")
		{
			public.GET("/article/page", articleHandler.PublicPage)        // 文章分页（仅已发布）
			public.GET("/article/dates", articleHandler.PublicDates)      // 发布日历（按年月）
			public.GET("/article/es/search", articleHandler.PublicSearch) // ES全文搜索（仅已发布）
			public.GET("/article/:id", articleHandler.PublicDetail)       // 文章详情
			public.GET("/category/tree", categoryHandler.Tree)            // 分类树
			public.GET("/site/config", siteHandler.PublicConfig)          // 站点配置
			public.GET("/profile", siteProfileHandler.PublicProfile)      // 站长个人资料
			public.GET("/links", siteLinkHandler.PublicLinks)             // 常用网站
		}

		// ==================== 需要登录的接口 ====================
		authorized := api.Group("")
		authorized.Use(middleware.JWTAuth(redisClient))
		// RBAC 中间件：markPerm 先按路由声明表写入 perm 标识，RBACAuth 再校验；
		// 超管（JWT isSuper=1）与未标注路由（perm 为空）均直接放行，读接口不受影响
		authorized.Use(markPerm())
		authorized.Use(middleware.RBACAuth(redisClient, roleRepo))
		{
			// 认证
			authorized.POST("/auth/logout", authHandler.Logout)
			authorized.GET("/auth/userinfo", authHandler.GetUserInfo)

			// 菜单
			menu := authorized.Group("/menu")
			{
				menu.GET("/list", menuHandler.List)
				menu.POST("/create", menuHandler.Create)
				menu.PUT("/update", menuHandler.Update)
				menu.DELETE("/:id", menuHandler.Delete)
			}

			// 角色管理
			role := authorized.Group("/role")
			{
				role.GET("/list", roleHandler.List)
				role.POST("/create", roleHandler.Create)
				role.PUT("/update", roleHandler.Update)
				role.DELETE("/:id", roleHandler.Delete)
			}

			// 用户管理
			user := authorized.Group("/user")
			{
				user.GET("/page", userHandler.Page)
				user.POST("/create", userHandler.Create)
				user.PUT("/update", userHandler.Update)
				user.PUT("/resetPwd/:id", userHandler.ResetPassword)
				user.PUT("/status/:id", userHandler.UpdateStatus)
				user.DELETE("/:id", userHandler.Delete)
			}

			// 分类管理
			category := authorized.Group("/category")
			{
				category.GET("/tree", categoryHandler.Tree)
				category.POST("/create", categoryHandler.Create)
				category.PUT("/update", categoryHandler.Update)
				category.DELETE("/:id", categoryHandler.Delete)
			}

			// 文章管理
			article := authorized.Group("/article")
			{
				article.GET("/page", articleHandler.Page)
				article.GET("/es/search", articleHandler.SearchArticle)
				article.POST("/create", articleHandler.Create)
				article.PUT("/update", articleHandler.Update)
				article.DELETE("/:id", articleHandler.Delete)
				article.POST("/es/rebuild", articleHandler.RebuildIndex)
			}

			// 站点配置
			site := authorized.Group("/site")
			{
				site.GET("/config", siteHandler.GetConfig)
				site.PUT("/config", siteHandler.SaveConfig)
				site.GET("/stats", siteHandler.Stats)
				site.POST("/testOss", siteHandler.TestOss)
				site.GET("/profile", siteProfileHandler.GetProfile)
				site.PUT("/profile", siteProfileHandler.SaveProfile)

				// 常用网站（友情链接）
				site.GET("/links", siteLinkHandler.List)
				site.POST("/links", siteLinkHandler.Create)
				site.PUT("/links/:id", siteLinkHandler.Update)
				site.DELETE("/links/:id", siteLinkHandler.Delete)
			}

			// 个人信息
			authorized.GET("/profile/info", profileHandler.GetInfo)
			authorized.PUT("/profile/update", profileHandler.Update)

			// 系统运行信息（仪表盘）
			authorized.GET("/system/overview", systemHandler.Overview)

			// 文件上传
			authorized.POST("/upload/oss", uploadHandler.UploadOSS)
		}
	}

	// Swagger 文档路由（仅在配置启用时）
	if cfg.Swagger.Enabled {
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	return router
}
