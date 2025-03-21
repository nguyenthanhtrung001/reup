package httpserver

import (
	"book-store/pkg/jwt"
	"context"
	"net/http"

	bookHTTP "book-store/internal/book/delivery/http"
	userHTTP "book-store/internal/user/delivery/http"

	bookRepo "book-store/internal/book/repository/mongo"
	userRepo "book-store/internal/user/repository/mongo"

	bookUseCase "book-store/internal/book/usecase"
	userUseCase "book-store/internal/user/usecase"

	"book-store/internal/middleware"

	bookProd "book-store/internal/book/delivery/rabbitmq/producer"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	// _ "book-store/docs"
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
