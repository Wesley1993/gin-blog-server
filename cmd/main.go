package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"wuzhispace.com/config"
	"wuzhispace.com/internal/handler"
	"wuzhispace.com/internal/repository"
	"wuzhispace.com/internal/service"
	"wuzhispace.com/pkg/elasticsearch"
	"wuzhispace.com/pkg/jwt"
	"wuzhispace.com/pkg/logger"
	"wuzhispace.com/pkg/migration"
	pkgredis "wuzhispace.com/pkg/redis"
	"wuzhispace.com/router"

	_ "wuzhispace.com/docs"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// @title Blog Admin API
// @version 1.0
// @description 博客后台管理系统 API 文档
// @host localhost:3000
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description 输入 Bearer {token}
func main() {
	// 记录启动时间，用于计算运行时长（仪表盘展示）
	startTime := time.Now()

	// 1. 初始化配置
	cfg, err := config.InitConfig()
	if err != nil {
		fmt.Printf("初始化配置失败: %v\n", err)
		os.Exit(1)
	}

	// 1.5 初始化日志
	logger.InitLogger(cfg.Log.Level)

	// 1.6 初始化 JWT（密钥与过期时间）
	jwt.Init(cfg.JWT.Secret, cfg.JWT.Expire)

	// 2. 初始化 PostgreSQL
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User,
		cfg.Database.Password, cfg.Database.DBName)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Fatal("连接数据库失败", "error", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		logger.Fatal("获取数据库连接池失败", "error", err)
	}
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.Database.ConnMaxLifetime) * time.Second)
	logger.Info("PostgreSQL 连接成功")

	// 2.5 执行数据库迁移（失败不阻止启动，仅记录错误）
	if err := migration.Run(sqlDB, cfg.Migration.Dir); err != nil {
		logger.Errorf("数据库迁移失败: %v", err)
	}

	// 3. 初始化 Redis
	pkgredis.InitRedis(cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.DB, cfg.Redis.PoolSize, cfg.Redis.MinIdleConns)
	redisClient := pkgredis.GetClient()
	logger.Info("Redis 连接成功")

	// 4. 初始化 ES（可选，失败不阻止启动）
	var esClient *elasticsearch.ESClient
	esClient, err = elasticsearch.NewESClient(&cfg.Elasticsearch)
	if err != nil {
		logger.Warn("ES初始化失败，搜索将降级到PG", "error", err)
		esClient = nil
	} else {
		esClient.CheckAndInit(context.Background())
		if !esClient.Available {
			logger.Warn("ES不可用，搜索将降级到PG")
		}
	}

	// 5. 初始化 Repository 层
	userRepo := repository.NewUserRepository(db)
	roleRepo := repository.NewRoleRepository(db)
	menuRepo := repository.NewMenuRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	articleRepo := repository.NewArticleRepository(db)
	siteConfigRepo := repository.NewSiteConfigRepository(db)
	siteProfileRepo := repository.NewSiteProfileRepository(db)
	siteLinkRepo := repository.NewSiteLinkRepository(db)

	// 6. 初始化 Service 层
	authService := service.NewAuthService(userRepo, roleRepo, menuRepo, db, redisClient)
	userService := service.NewUserService(userRepo, roleRepo, db, redisClient)
	roleService := service.NewRoleService(roleRepo, userRepo, db, redisClient)
	menuService := service.NewMenuService(menuRepo, roleRepo, userRepo, redisClient)
	categoryService := service.NewCategoryService(categoryRepo, db)
	articleService := service.NewArticleService(articleRepo, categoryRepo, db, esClient)
	siteService := service.NewSiteConfigService(siteConfigRepo, db, redisClient)
	siteProfileService := service.NewSiteProfileService(siteProfileRepo, db, redisClient)
	siteLinkService := service.NewSiteLinkService(siteLinkRepo)
	profileService := service.NewProfileService(userRepo, db, redisClient)
	systemService := service.NewSystemService(startTime, db, redisClient, esClient)

	// 7. 初始化 Handler 层
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)
	roleHandler := handler.NewRoleHandler(roleService)
	menuHandler := handler.NewMenuHandler(menuService)
	categoryHandler := handler.NewCategoryHandler(categoryService)
	articleHandler := handler.NewArticleHandler(articleService)
	siteHandler := handler.NewSiteHandler(siteService)
	siteProfileHandler := handler.NewSiteProfileHandler(siteProfileService)
	siteLinkHandler := handler.NewSiteLinkHandler(siteLinkService)
	uploadHandler := handler.NewUploadHandler(siteConfigRepo)
	profileHandler := handler.NewProfileHandler(profileService)
	systemHandler := handler.NewSystemServiceHandler(systemService)

	// 8. 设置路由
	r := router.SetupRouter(
		authHandler, userHandler, roleHandler, menuHandler,
		categoryHandler, articleHandler, siteHandler, siteProfileHandler, siteLinkHandler, uploadHandler, profileHandler, systemHandler,
		roleRepo, cfg,
	)

	// 9. 启动 HTTP 服务（支持优雅退出）
	addr := fmt.Sprintf(":%d", cfg.App.Port)
	logger.Infof("%s v%s 启动在 %s", cfg.App.Name, cfg.App.Version, addr)

	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("启动服务失败", "error", err)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("正在关闭服务器...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Errorf("服务器强制关闭: %v", err)
	}

	// 关闭 DB 连接
	if sqlDB, err := db.DB(); err == nil {
		_ = sqlDB.Close()
	}
	// 关闭 Redis 连接
	if redisClient != nil {
		_ = redisClient.Close()
	}
	logger.Info("服务器已退出")
}
