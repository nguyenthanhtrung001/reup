package httpserver

import (
	"context"
	"net/http"
	"reup/pkg/jwt"
	"reup/pkg/telegram"

	bookHTTP "reup/internal/book/delivery/http"
	userHTTP "reup/internal/user/delivery/http"

	bookRepo "reup/internal/book/repository/mongo"
	userRepo "reup/internal/user/repository/mongo"

	bookUseCase "reup/internal/book/usecase"
	userUseCase "reup/internal/user/usecase"

	"reup/internal/middleware"

	bookProd "reup/internal/book/delivery/rabbitmq/producer"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	// _ "reup/docs"
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
	bookRepository := bookRepo.New(srv.l, srv.database, jwtManager)
	if bookRepository == nil {
		srv.l.Fatal(context.Background(), "Failed to initialize book repository")
	}

	userRepository := userRepo.New(srv.l, srv.database, jwtManager)
	if userRepository == nil {
		srv.l.Fatal(context.Background(), "Failed to initialize user repository")
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
	bookProd := bookProd.New(srv.l, srv.amqpConn)
	if err := bookProd.Run(); err != nil {
		srv.l.Fatal(context.Background(), "Failed to run book producer: %v", err)
	}

	// UseCases
	bookUC := bookUseCase.New(srv.l, bookRepository, encrypter, bookProd)
	if bookUC == nil {
		srv.l.Fatal(context.Background(), "Failed to initialize book use case")
	}

	userUC := userUseCase.New(srv.l, userRepository, encrypter)
	if userUC == nil {
		srv.l.Fatal(context.Background(), "Failed to initialize user use case")
	}

	// Handlers
	bookH := bookHTTP.New(srv.l, bookUC)
	if bookH == nil {
		srv.l.Fatal(context.Background(), "Failed to initialize book handler")
	}

	userH := userHTTP.New(srv.l, userUC)
	if userH == nil {
		srv.l.Fatal(context.Background(), "Failed to initialize user handler")
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

	// Map book routes
	bookHTTP.MapbookRoutes(api.Group("/book"), bookH, mw)
	userHTTP.MapUserRoutes(api.Group("/user"), userH, mw)
}
