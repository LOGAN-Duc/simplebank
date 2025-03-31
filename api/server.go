package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/techschool/simplebank/db/sqlc"
)

type Server struct {
	// giúp kết nối csdl khi sử lý các yêu cầu API
	store *db.Store
	// gửi từng yêu cầu API đến trình xử lý thích hợp
	router *gin.Engine
}

// thiết lập các tuyến API HTTP
func NewServer(store *db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	//thêm route API
	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccount)
	router.PUT("/accounts/:id", server.updateAccount)
	router.DELETE("/accounts/:id", server.deleteAccount)

	server.router = router
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error :": err.Error()}
}
