package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"github.com/techschool/simplebank/Utill"
	db "github.com/techschool/simplebank/db/sqlc"
	"net/http"
)

type createUserRequest struct {
	Username       string `json:"username" binding:"required,alphanum"`
	HashedPassword string `json:"password" binding:"required,min=6"`
	FullName       string `json:"full_name" binding:"required"`
	Email          string `json:"email" binding:"required,email"`
}

func (server *Server) createUser(ctx *gin.Context) {
	var rep createUserRequest
	if err := ctx.ShouldBindJSON(&rep); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	hashPassword, err := Utill.HashPassword(rep.HashedPassword)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	arg := db.CreateUserParams{
		Username:       rep.Username,
		HashedPassword: hashPassword,
		FullName:       rep.FullName,
		Email:          rep.Email,
	}
	//arg = db.CreateUserParams{}
	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, user)

}
