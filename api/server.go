package api

import (
	"log"
	db "simplebank/db/sqlc"
	"simplebank/token"
	"simplebank/util"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Server serves HTTP requests for our banking service
type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

// NewServer creates a new HTTP server and set up routing
func NewServer(config util.Config, store db.Store) (*Server, error) {
	server := &Server{
		config: config,
		store:  store,
	}

	server.setupTokenMaker()
	server.setupRouter()
	server.setupValidator()
	return server, nil
}

// Start runs the HTTP server on a specific address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func (server *Server) setupRouter() {
	router := gin.Default()

	// 不需要 Auth 中间件,任何用户都可以注册和登录
	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)

	// 对需要 Auth 中间件的路由进行分组
	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))
	authRoutes.POST("/accounts", server.createAccount)
	authRoutes.GET("/accounts/:id", server.getAccount)
	authRoutes.GET("/accounts", server.listAccount)
	authRoutes.POST("/transfers", server.createTransfer)

	server.router = router
}

// 注册验证器
func (server *Server) setupValidator() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		if err := v.RegisterValidation("currency", validCurrency); err != nil {
			log.Fatalf("failed to register currency validator: %v", err)
		}
	}
}

// 设置Token发放
func (server *Server) setupTokenMaker() {
	tokenMaker, err := token.NewPasetoMaker(server.config.TokenSymmetricKey)
	if err != nil {
		log.Fatal("cannot create token maker: %w", err)
	}
	server.tokenMaker = tokenMaker
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
