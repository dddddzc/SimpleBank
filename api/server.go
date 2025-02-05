package api

import (
	"log"
	db "simplebank/db/sqlc"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Server serves HTTP requests for our banking service
type Server struct {
	store  db.Store
	router *gin.Engine
}

// NewServer creates a new HTTP server and set up routing
func NewServer(store db.Store) *Server {
	server := &Server{
		store:  store,
		router: gin.Default(),
	}
	server.setupRouter()
	server.setupValidator()
	return server
}

// Start runs the HTTP server on a specific address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func (server *Server) setupRouter() {
	// 从请求体(Request Body)中解析 JSON 数据
	server.router.POST("/accounts", server.createAccount)
	// 从 URL 路径中提取参数
	server.router.GET("/accounts/:id", server.getAccount)
	// 从 URL 查询字符串中提取参数
	server.router.GET("/accounts", server.listAccount)
	server.router.POST("/transfers", server.createTransfer)
	server.router.POST("/users", server.createUser)
}

// 注册验证器
func (server *Server) setupValidator() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		if err := v.RegisterValidation("currency", validCurrency); err != nil {
			log.Fatalf("failed to register currency validator: %v", err)
		}
	}
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
