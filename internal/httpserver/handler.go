package httpserver

import (
	"context"
	"net/http"

	"github.com/nguyenthanhtrung001/reup/pkg/jwt"
	"github.com/nguyenthanhtrung001/reup/pkg/telegram"

	scanVideoHTTP "github.com/nguyenthanhtrung001/reup/internal/scan_video/delivery/http"
	userHTTP "github.com/nguyenthanhtrung001/reup/internal/user/delivery/http"

	scanVideoRepo "github.com/nguyenthanhtrung001/reup/internal/scan_video/repository/mongo"
	userRepo "github.com/nguyenthanhtrung001/reup/internal/user/repository/mongo"

	scanVideoUseCase "github.com/nguyenthanhtrung001/reup/internal/scan_video/usecase"
	userUseCase "github.com/nguyenthanhtrung001/reup/internal/user/usecase"

	"github.com/nguyenthanhtrung001/reup/internal/middleware"

	scanVideoProd "github.com/nguyenthanhtrung001/reup/internal/scan_video/delivery/rabbitmq/producer"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	// _ "github.com/nguyenthanhtrung001/reup/docs"
)

func (srv HTTPServer) mapHandlers() {
	// Swagger
	srv.gin.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// JWT manager
	if srv.jwtSecretKey == "" {
		srv.l.Fatal(context.Background(), "JWT secret key is not defined")
	}
	jwtManager, err := jwt.NewJWTMaker(srv.jwtSecretKey)
	if err != nil {
		srv.l.Fatal(context.Background(), "Failed to create JWT manager: %v", err)
	}

	// Encrypter
	encrypter := srv.encrypter
	if encrypter == nil {
		srv.l.Fatal(context.Background(), "Encrypter is not defined")
	}

	// Repositories

	userRepository := userRepo.New(srv.l, srv.database, jwtManager)
	if userRepository == nil {
		srv.l.Fatal(context.Background(), "Failed to initialize user repository")
	}
	scanVideoRepository := scanVideoRepo.New(srv.l, srv.database, jwtManager)
	if scanVideoRepository == nil {
		srv.l.Fatal(context.Background(), "Failed to initialize scanVideoRepository repository")
	}

	chatIDs := telegram.ChatIDs{
		ReportBug:     srv.telegram.ChatIDs.ReportBug,
		ReportPayment: srv.telegram.ChatIDs.ReportPayment,
	}
	telegram := telegram.New(srv.telegram.BotKey, chatIDs)
	srv.gin.Use(middleware.Recovery(telegram, srv.telegram.ChatIDs.ReportBug))
	// srv.l.Fatal(context.Background(), telegram)

	// Producer
	if srv.amqpConn == nil {
		srv.l.Fatal(context.Background(), "AMQP connection is nil")
	}

	scanVideoProd := scanVideoProd.New(srv.l, srv.amqpConn)
	if err := scanVideoProd.Run(); err != nil {
		srv.l.Fatal(context.Background(), "Failed to run sacnVideoProd producer: %v", err)
	}

	// UseCases
	userUC := userUseCase.New(srv.l, userRepository, encrypter)
	if userUC == nil {
		srv.l.Fatal(context.Background(), "Failed to initialize user use case")
	}
	scanVideoUC := scanVideoUseCase.New(srv.l, scanVideoRepository, encrypter, scanVideoProd, telegram, scanVideoUseCase.TeleChat{
		NotifiChatID: srv.telegram.ReportPayment,
		GroupChat1:   srv.telegram.GroupChat1,
		GroupChat2:   srv.telegram.GroupChat2,
		GroupChat3:   srv.telegram.GroupChat3,
	})
	if scanVideoUC == nil {
		srv.l.Fatal(context.Background(), "Failed to initialize scanVideoUC use case")
	}

	// Handlers

	userH := userHTTP.New(srv.l, userUC)
	if userH == nil {
		srv.l.Fatal(context.Background(), "Failed to initialize user handler")
	}

	scanVideoH := scanVideoHTTP.New(srv.l, scanVideoUC)
	if scanVideoH == nil {
		srv.l.Fatal(context.Background(), "Failed to initialize scanVideoH handler")
	}

	// API group
	api := srv.gin.Group("/api/v1")

	// Add middleware to handle OPTIONS requests
	api.OPTIONS("/*path", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization, Sec-Key, Accept, Referer, User-Agent")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.JSON(http.StatusOK, gin.H{})
	})

	// Middlewares
	mw := middleware.New(srv.l, jwtManager, userUC, encrypter, srv.secretKey)

	// CORS
	api.Use(mw.Cors())

	// System maintenance
	api.Use(mw.SystemMaintenance())

	// Map  routes
	userHTTP.MapUserRoutes(api.Group("/user"), userH, mw)
	scanVideoHTTP.MapScanVideoRoutes(api.Group("/scan"), scanVideoH, mw)
	scanVideoHTTP.MapProxyRoutes(api.Group("/"), scanVideoH, mw)
}
