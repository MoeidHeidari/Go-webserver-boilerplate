package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"main/bootstrap"
)

func main() {
	_ = godotenv.Load()
	err := bootstrap.RootApp.Execute()
	if err != nil {
		return
	}
}
