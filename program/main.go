package main

import "go-gin-sqlserver/program/application"

func main() {
	application.Route()
	application.Listen()
}
