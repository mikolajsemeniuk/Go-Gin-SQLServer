package controllers

import (
	"go-gin-sqlserver/program/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddUser(context *gin.Context) {
	database.Client.Ping()
	context.JSON(http.StatusOK, "hi")
}
