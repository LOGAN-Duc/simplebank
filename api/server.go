package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/techschool/simplebank/db/sqlc"
	"github.com/techschool/simplebank/token"
	"github.com/techschool/simplebank/util"
)

type Server struct {
	// giúp kết nối csdl khi sử lý các yêu cầu API
	store db.Store
	// gửi từng yêu cầu API đến trình xử lý thích hợp
	router *gin.Engine
	//tao khoa
	tokenMaker token.Maker
	config     util.Config
}

// thiết lập các tuyến API HTTP
func NewServer(config util.Config, store db.Store) (*Server, error) {

	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)

	if err != nil {
		return nil, fmt.Errorf("can not create token maker : %w", err)
	}
	server := &Server{store: store,
		tokenMaker: tokenMaker,
		config:     config}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}
	server.SetupRouter()
	return server, nil
}
func (server *Server) SetupRouter() {
	router := gin.Default()
	//add middleaware
	auThRouter := router.Group("/").Use(authMiddeleware(server.tokenMaker))
	//thêm route API
	//account
	auThRouter.POST("/accounts", server.createAccount)
	auThRouter.GET("/accounts/:id", server.getAccount)
	auThRouter.GET("/accounts", server.listAccount)
	auThRouter.PUT("/accounts/:id", server.updateAccount)
	auThRouter.DELETE("/accounts/:id", server.deleteAccount)
	//transfer
	auThRouter.POST("/transfers", server.createTransfer)
	//router.GET("/accounts/:id", server.getAccount)
	//router.GET("/accounts", server.listAccount)
	//router.PUT("/accounts/:id", server.updateAccount)
	//router.DELETE("/accounts/:id", server.deleteAccount)
	//user
	router.POST("/users", server.createUser)
	router.POST("/users/login", server.LoginUser)

	//router.GET("/accounts/:id", server.getAccount)
	//router.GET("/accounts", server.listAccount)
	//router.PUT("/accounts/:id", server.updateAccount)
	//router.DELETE("/accounts/:id", server.deleteAccount)
	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error :": err.Error()}
}
