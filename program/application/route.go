package application

import "go-gin-sqlserver/program/controllers"

func Route() {
	router.GET("/", controllers.AddUser)
}
