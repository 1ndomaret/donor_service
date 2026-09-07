package main

import (
	"github.com/labstack/echo/v5"
)

func main() {
	// db, err := config.InitDB()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// if err := db.AutoMigrate(); err != nil {
	// 	log.Fatal("failed to migrate database:", err)
	// }

	e := echo.New()

	// router.Register(e,

	// )
	if err := e.Start(":1323"); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
